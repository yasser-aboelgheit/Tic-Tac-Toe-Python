package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	serviceConfig "gitlab.com/startupbuilder/uda/internal/configs"
	"gitlab.com/startupbuilder/startupbuilder/pkg/config"
	logger "gitlab.com/startupbuilder/startupbuilder/pkg/logger"
)

func NewRootCmd() *cobra.Command {
	var (
		lgr         logger.Log
		envFilePath string
	)

	cfg := &serviceConfig.GeneralConfig{}

	cmd := &cobra.Command{
		Short: "Start service.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// read configurations.
			err := config.ReadConfig(cfg,
				config.WithFilePath(envFilePath),
				config.WithPrefix(serviceConfig.ENVFILEPREFIX),
			)
			if err != nil {
				return fmt.Errorf("could not read auth service configs: %w", err)
			}

			// config logger.
			lgr = logger.NewLogger(
				cfg.Logger,
				logger.WithEnvironment(cfg.Service.Environment),
				logger.WithVersion(cfg.Service.Version),
				logger.WithServiceName(cfg.Service.Name),
			)

			if cfg.Logger.PrettyPrint {
				lgr = lgr.WithPrettyOutput()
				lgr.Debug(cmd.Context(), "prettify")
			}

			return nil
		},
	}

	cmd.AddCommand(newHTTPCmd(cfg, &lgr))
	cmd.PersistentFlags().StringVarP(&envFilePath, "env", "e", "", "env file path")

	return cmd
}

