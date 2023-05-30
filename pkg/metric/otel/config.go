package otel

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// config contains metrics configuration.
type config struct {
	exportInterval           time.Duration
	serviceName              string
	environment              string
	version                  string
	useCumulativeTemporality bool
	attrs                    []attribute.KeyValue
}

// Option is a metrics configuration option.
type Option interface {
	apply(*config) error
}

// optionFunc applies options to a config.
type optionFunc func(*config) error

// apply applies options to a config.
func (o optionFunc) apply(cfg *config) error {
	return o(cfg)
}

// WithExportInterval returns an Option that sets metrics export interval.
func WithExportInterval(interval time.Duration) Option { //nolint:ireturn
	return optionFunc(func(cfg *config) error {
		if interval > 0 {
			cfg.exportInterval = interval
		}

		return nil
	})
}

// WithServiceNameAttribute returns an Option that registers a global Service
// Name attribute that will be attached to all metrics.
func WithServiceNameAttribute(serviceName string) Option { //nolint:ireturn
	return optionFunc(func(cfg *config) error {
		cfg.serviceName = serviceName

		return nil
	})
}

// WithEnvironmentAttribute returns an Option that registers a global
// Environment attribute that will be attached to all metrics.
func WithEnvironmentAttribute(environment string) Option { //nolint:ireturn
	return optionFunc(func(cfg *config) error {
		cfg.environment = environment

		return nil
	})
}

// WithVersionAttribute returns an Option that registers a global Version
// attribute that will be attached to all metrics.
func WithVersionAttribute(version string) Option { //nolint:ireturn
	return optionFunc(func(cfg *config) error {
		cfg.version = version

		return nil
	})
}

// WithCumulativeTemporality returns an Option that sets aggregation
// temporality to Cumulative.
func WithCumulativeTemporality() Option { //nolint:ireturn
	return optionFunc(func(cfg *config) error {
		cfg.useCumulativeTemporality = true

		return nil
	})
}

// WithAttributes adds attributes to metrics.
func WithAttributes(attrs ...attribute.KeyValue) Option {
	return optionFunc(func(cfg *config) error {
		cfg.attrs = append(cfg.attrs, attrs...)

		return nil
	})
}
