package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/startupbuilder/uda/internal/configs"
	"gitlab.com/startupbuilder/uda/internal/server"
	httpLib "gitlab.com/startupbuilder/startupbuilder/pkg/http"
	"gitlab.com/startupbuilder/startupbuilder/pkg/http/pprof"
	"gitlab.com/startupbuilder/startupbuilder/pkg/logger"
	"gitlab.com/startupbuilder/startupbuilder/pkg/tracer"
	"golang.org/x/sync/errgroup"
)

func newHTTPCmd(cfg *configs.AuthConfig, lgr *logger.Log) *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Serving HTTP API.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := logger.ContextWithAttributes(
				cmd.Context(),
				map[string]interface{}{
					"command": "http",
				},
			)

			lgr.Infow(ctx, "Starting Server", map[string]interface{}{
				"tracer.endpoint": cfg.Tracer.ReceiverEndpoint,
				"http.port":       cfg.HTTP.Port,
				"http.host":       cfg.HTTP.Host,
			})

			tracerProvider, err := tracer.NewTraceProvider(
				ctx,
				tracer.WithReceiverEndpoint(cfg.Tracer.ReceiverEndpoint),
				tracer.WithServiceName(cfg.Service.Name),
				tracer.WithEnvironment(cfg.Service.Environment),
				tracer.WithLogger(lgr),
				tracer.WithVersion(cfg.Service.Version),
			)
			if err != nil {
				return fmt.Errorf("failed to setup tracer: %w", err)
			}

			defer tracerProvider.Stop()

			mux := server.NewHTTPHandlers(lgr, tracerProvider.Tracer("http"))

			group, ctx := errgroup.WithContext(ctx)

			pprofServ := pprof.NewPPROFServer(cfg.Pprof, lgr)
			group.Go(pprofServ.Start)
			group.Go(pprofServ.SignalShutdown)

			serv := httpLib.NewHTTPServer(
				mux,
				lgr,
				cfg.HTTP.Options()...,
			)
			group.Go(serv.Start)
			group.Go(serv.SignalShutdown)

			if err := group.Wait(); err != nil {
				lgr.Error(ctx, err, "server error")
			}

			lgr.Info(ctx, "Shutting Down")

			return nil
		},
	}
}
