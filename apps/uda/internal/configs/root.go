package configs

import (
	"gitlab.com/startupbuilder/startupbuilder/pkg/config"
	httpConfig "gitlab.com/startupbuilder/startupbuilder/pkg/http"
	pprofConfig "gitlab.com/startupbuilder/startupbuilder/pkg/http/pprof"
	loggerConfig "gitlab.com/startupbuilder/startupbuilder/pkg/logger"
	tracerConfig "gitlab.com/startupbuilder/startupbuilder/pkg/tracer"
)

const ENVFILEPREFIX = "AUTH"

type AuthConfig struct {
	HTTP    *httpConfig.Config    `mapstructure:"http"`
	Pprof   *pprofConfig.Config   `mapstructure:"pprof"`
	Logger  *loggerConfig.Config  `mapstructure:"logger"`
	Service *config.ServiceConfig `mapstructure:"service"`
	Tracer  *tracerConfig.Config  `mapstructure:"tracer"`
}

func (cfg *AuthConfig) Defaults() {
	cfg.Service = config.NewServiceConfig()
	cfg.Service.Name = "auth"

	cfg.Logger = loggerConfig.NewConfig()

	cfg.HTTP = httpConfig.NewConfig()

	cfg.Pprof = pprofConfig.NewConfig()

	cfg.Tracer = &tracerConfig.Config{}
	cfg.Tracer.Defaults()
}