// Basic logger configuration.
package logger

import "fmt"

type Config struct {
	Level       string `mapstructure:"LEVEL"`
	PrettyPrint bool   `mapstructure:"PRETTY0PRINT"`
	environemnt string
	serviceName string
	version     string
}

func NewConfig() *Config {
	return &Config{
		Level:       string(DebugLevel),
		environemnt: "un-specified",
		serviceName: "un-specified",
		version:     "un-specified",
	}
}

func (cfg *Config) apply(opts ...Option) error {
	for _, opt := range opts {
		err := opt(cfg)
		if err != nil {
			return fmt.Errorf("can not apply settings: %w", err)
		}
	}

	return nil
}

type Option func(*Config) error

func WithEnvironment(env string) Option {
	return func(c *Config) error {
		if env != "" {
			c.environemnt = env
		}

		return nil
	}
}

func WithServiceName(name string) Option {
	return func(c *Config) error {
		if name != "" {
			c.serviceName = name
		}

		return nil
	}
}

func WithVersion(version string) Option {
	return func(c *Config) error {
		if version != "" {
			c.version = version
		}

		return nil
	}
}
