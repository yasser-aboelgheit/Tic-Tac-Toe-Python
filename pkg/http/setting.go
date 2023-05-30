package http

import (
	"fmt"
	"time"
)

type Config struct {
	Host               string        `mapstructure:"HOST"`
	Port               int           `mapstructure:"PORT"`
	ReadTimeout        time.Duration `mapstructure:"READ0TIMEOUT"`
	ReadHeaderTimeout  time.Duration `mapstructure:"READ0HEADER0TIMEOUT"`
	WriteTimeout       time.Duration `mapstructure:"WRITE0TIMEOUT"`
	MaxShutdownTimeout time.Duration `mapstructure:"MAX0SHUTDOWN0TIMEOUT"`
	BasePrefix         string        `mapstructure:"BASE0URL"`
}

func NewConfig() *Config {
	return &Config{
		Port:               8000,
		ReadHeaderTimeout:  time.Millisecond,
		ReadTimeout:        time.Millisecond * 10,
		WriteTimeout:       time.Millisecond * 40,
		MaxShutdownTimeout: time.Second * 1,
	}
}

func (s *Config) apply(opts ...Option) error {
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return fmt.Errorf("settings error on HTTP: %w", err)
		}
	}

	return nil
}

type Option func(*Config) error

func WithHost(host string) Option {
	return func(s *Config) error {
		s.Host = host

		return nil
	}
}

func WithPort(port int) Option {
	return func(s *Config) error {
		s.Port = port

		return nil
	}
}

func WithReadTimeout(input time.Duration) Option {
	return func(s *Config) error {
		s.ReadTimeout = input

		return nil
	}
}

func WithReadHeaderTimeout(input time.Duration) Option {
	return func(s *Config) error {
		s.ReadHeaderTimeout = input

		return nil
	}
}

func WithWriteTimeout(input time.Duration) Option {
	return func(s *Config) error {
		s.WriteTimeout = input

		return nil
	}
}

func WithShutdownTimeout(input time.Duration) Option {
	return func(s *Config) error {
		s.MaxShutdownTimeout = input

		return nil
	}
}

func WithBaseURL(baseUrl string) Option {
	return func(s *Config) error {
		s.BasePrefix = baseUrl

		return nil
	}
}
