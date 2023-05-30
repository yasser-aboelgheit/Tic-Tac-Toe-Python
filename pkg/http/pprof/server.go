package pprof

import (
	"fmt"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	gorillaMux "github.com/gorilla/mux"
	httpProvider "gitlab.com/startupbuilder/startupbuilder/pkg/http"
)

type PPROFServer struct {
	cfg   *Config
	pprof *httpProvider.HTTPServer
	lgr   handlers.RecoveryHandlerLogger
}

func NewPPROFServer(cfg *Config, lgr handlers.RecoveryHandlerLogger) *PPROFServer {
	return &PPROFServer{
		cfg: cfg,
		lgr: lgr,
	}
}

func (server *PPROFServer) Start() func() error {
	mux := gorillaMux.NewRouter()

	mux.HandleFunc("/debug/pprof", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)

	server.pprof = httpProvider.NewHTTPServer(
		server.lgr,
		httpProvider.WithHost(server.cfg.HTTP.Host),
		httpProvider.WithPort(server.cfg.HTTP.Port),
		httpProvider.WithReadHeaderTimeout(server.cfg.HTTP.ReadHeaderTimeout),
		httpProvider.WithReadTimeout(server.cfg.HTTP.ReadTimeout),
		httpProvider.WithWriteTimeout(server.cfg.HTTP.WriteTimeout),
		httpProvider.WithShutdownTimeout(server.cfg.HTTP.MaxShutdownTimeout),
	)

	return server.pprof.Start(mux)
}

func (server *PPROFServer) Close() error {
	return server.pprof.Close()
}

func (server *PPROFServer) SignalShutdown() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	if err := server.pprof.Close(); err != nil {
		return fmt.Errorf("failed to graceful shutdown on OS signal: %w", err)
	}

	return nil
}
