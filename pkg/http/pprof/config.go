package pprof

import (
	"time"

	"gitlab.com/startupbuilder/startupbuilder/pkg/http"
)

type Config struct {
	HTTP *http.Config `mapstructure:",squash"`
}

func NewConfig() *Config {
	cfg := &Config{
		HTTP: http.NewConfig(),
	}

	cfg.HTTP.Host = "localhost"
	cfg.HTTP.Port = 3002
	cfg.HTTP.WriteTimeout = time.Minute * 10 // enable most of pprof writer timeouts

	return cfg
}

