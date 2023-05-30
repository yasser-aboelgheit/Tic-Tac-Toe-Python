package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type HTTPServer struct {
	stngs *Config
	http  *http.Server
	lgr   handlers.RecoveryHandlerLogger
}

// NewHTTPServer will create a default settings HTTPServer according to our standard.
// to change any of the defaults pass it as option.
func NewHTTPServer(
	lgr handlers.RecoveryHandlerLogger,
	opts ...Option,
) *HTTPServer {
	stngs := NewConfig()
	if err := stngs.apply(opts...); err != nil {
		panic(err)
	}

	return &HTTPServer{
		stngs: stngs,
		lgr:   lgr,
	}
}

func (server *HTTPServer) Start(mux *mux.Router) func() error {
	newMux := http.NewServeMux()
	newMux.Handle(server.stngs.BasePrefix+"/", http.StripPrefix(server.stngs.BasePrefix, mux))

	server.http = &http.Server{
		Handler: handlers.RecoveryHandler(
			handlers.RecoveryLogger(server.lgr),
			handlers.PrintRecoveryStack(true),
		)(newMux),
		Addr:              fmt.Sprintf("%s:%d", server.stngs.Host, server.stngs.Port),
		ReadTimeout:       server.stngs.ReadTimeout,
		ReadHeaderTimeout: server.stngs.ReadHeaderTimeout,
		WriteTimeout:      server.stngs.WriteTimeout,
	}
	return func() error {
		err := server.http.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed starting http: %w", err)
		}

		return nil
	}
}

func (server *HTTPServer) Close() error {
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		server.stngs.MaxShutdownTimeout,
	)
	defer cancel()

	if err := server.http.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("HTTP shutdown error: %v", err)
	}

	return nil
}

func (server *HTTPServer) SignalShutdown() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	if err := server.Close(); err != nil {
		return fmt.Errorf("failed to graceful shutdown on OS signal: %w", err)
	}

	return nil
}
