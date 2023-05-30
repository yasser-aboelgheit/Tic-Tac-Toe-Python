package http

import (
	"fmt"
	"time"
)

type settings struct {
	Host               string
	Port               int
	ReadTimeout        time.Duration
	ReadHeaderTimeout  time.Duration
	WriteTimeout       time.Duration
	MaxShutdownTimeout time.Duration
}

func newSettings() *settings {
	return &settings{
		Host:               "127.0.0.1",
		Port:               8000,
		ReadHeaderTimeout:  time.Millisecond,
		ReadTimeout:        time.Millisecond * 10,
		WriteTimeout:       time.Millisecond * 40,
		MaxShutdownTimeout: time.Second * 1,
	}
}

func (s *settings) apply(opts ...Option) error {
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return fmt.Errorf("settings error on HTTP: %w", err)
		}
	}

	return nil
}

type Option func(*settings) error

func WithHost(host string) Option {
	return func(s *settings) error {
		s.Host = host

		return nil
	}
}

func WithPort(port int) Option {
	return func(s *settings) error {
		s.Port = port

		return nil
	}
}

func WithReadTimeout(input time.Duration) Option {
	return func(s *settings) error {
		s.ReadTimeout = input

		return nil
	}
}

func WithReadHeaderTimeout(input time.Duration) Option {
	return func(s *settings) error {
		s.ReadHeaderTimeout = input

		return nil
	}
}

func WithWriteTimeout(input time.Duration) Option {
	return func(s *settings) error {
		s.WriteTimeout = input

		return nil
	}
}

func WithShutdownTimeout(input time.Duration) Option {
	return func(s *settings) error {
		s.MaxShutdownTimeout = input

		return nil
	}
}
