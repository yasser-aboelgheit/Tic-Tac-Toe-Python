package http

import "time"

type Config struct {
	Host               string        `mapstructure:"HOST"`
	Port               int           `mapstructure:"PORT"`
	ReadTimeout        time.Duration `mapstructure:"READ0TIMEOUT"`
	ReadHeaderTimeout  time.Duration `mapstructure:"READ0HEADER0TIMEOUT"`
	WriteTimeout       time.Duration `mapstructure:"WRITE0TIMEOUT"`
	MaxShutdownTimeout time.Duration `mapstructure:"MAX0SHUTDOWN0TIMEOUT"`
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg Config) Options() []Option {
	opts := make([]Option, 5)

	if cfg.Host != "" {
		opts = append(opts, WithHost(cfg.Host))
	}

	if cfg.Port != 0 {
		opts = append(opts, WithPort(cfg.Port))
	}

	if cfg.ReadTimeout != 0 {
		opts = append(opts, WithReadTimeout(cfg.ReadTimeout))
	}

	if cfg.ReadHeaderTimeout != 0 {
		opts = append(opts, WithReadHeaderTimeout(cfg.ReadHeaderTimeout))
	}

	if cfg.WriteTimeout != 0 {
		opts = append(opts, WithWriteTimeout(cfg.WriteTimeout))
	}

	if cfg.MaxShutdownTimeout != 0 {
		opts = append(opts, WithShutdownTimeout(cfg.MaxShutdownTimeout))
	}

	return opts
}
