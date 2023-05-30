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

func (server *PPROFServer) Start() error {
	mux := gorillaMux.NewRouter()

	mux.HandleFunc("/debug/pprof", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)

	server.pprof = httpProvider.NewHTTPServer(mux, server.lgr, server.cfg.HTTP.Options()...)

	if err := server.pprof.Start(); err != nil {
		return fmt.Errorf("failure pprof: %w", err)
	}

	return nil
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
