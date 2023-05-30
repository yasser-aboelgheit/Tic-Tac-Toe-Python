package server

import (
	"fmt"
	"net/http"
	"time"

	gorillaMux "github.com/gorilla/mux"
	"gitlab.com/startupbuilder/uda/internal/endpoints"
	middlewreHTTP "gitlab.com/startupbuilder/startupbuilder/pkg/http/middleware"
	"gitlab.com/startupbuilder/startupbuilder/pkg/logger"
	"gitlab.com/startupbuilder/startupbuilder/pkg/tracer"
)

// NewHTTPHandlers create a paths routes
func NewHTTPHandlers(lgr logger.Logger, tracer tracer.Tracer, opts ...Option) *gorillaMux.Router {
	mux := gorillaMux.NewRouter()
	mux.Use(middlewreHTTP.MakeMuxLoggerMiddleware(lgr))
	mux.Use(middlewreHTTP.MakeJsonContentMiddleware(false))

	stngs := settings{}
	stngs.apply(mux, opts...)

	mux.HandleFunc("/health", endpoints.MakeHealthEndpoint(tracer)).Methods(http.MethodGet)

	return mux
}

func WithRequestTimeoutInSeconds(requestDuration int) Option {
	return func(s *settings, router *gorillaMux.Router) error {
		maxTimeoutDuration := time.Duration(requestDuration) * time.Second

		if maxTimeoutDuration > 0 {
			router.Use(middlewreHTTP.MakeTimeoutMiddleware(maxTimeoutDuration))
		}

		return nil
	}
}

type settings struct {
	maxTimeoutDuration time.Duration
}

type Option func(*settings, *gorillaMux.Router) error

func (stng *settings) apply(router *gorillaMux.Router, opts ...Option) error {
	for _, opt := range opts {
		err := opt(stng, router)
		if err != nil {
			return fmt.Errorf("can not apply handler option: %w", err)
		}
	}

	return nil
}
