package tracer

import (
	"time"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

const (
	defaultName              string        = "unknown"
	defaultEnvironment       string        = "local"
	defaultVersion           string        = "v0.0.0"
	defaultReceiver          string        = "localhost:4126"
	defaultSampleRatio       float64       = 1
	defaultSpanExportTimeout time.Duration = time.Millisecond * 100
	defaultResourceTimeout   time.Duration = time.Millisecond * 150
)

type settings struct {
	version         string
	name            string
	environment     string
	receiver        string
	sampleRatio     float64
	exporterTimeout time.Duration
	lgr             *errorLogger
	propagators     []propagation.TextMapPropagator
	resourceTimeout time.Duration
	traceInjectors  []AttributeInjector
	attributes      []attribute.KeyValue
}

func (settings *settings) Apply(opts ...ProviderOption) {
	for _, opt := range opts {
		opt(settings)
	}
}

func newDefaultSettings() settings {
	return settings{
		name:            defaultName,
		version:         defaultVersion,
		environment:     defaultEnvironment,
		sampleRatio:     defaultSampleRatio,
		receiver:        defaultReceiver,
		exporterTimeout: defaultSpanExportTimeout,
		resourceTimeout: defaultResourceTimeout,
		lgr:             &errorLogger{lgr: &nooplogger{}},
		propagators: []propagation.TextMapPropagator{
			propagation.TraceContext{},
			propagation.Baggage{},
		},
		traceInjectors: []AttributeInjector{},
		attributes:     []attribute.KeyValue{},
	}
}

type ProviderOption func(settings *settings) error

func WithEnvironment(environment string) ProviderOption {
	return func(settings *settings) error {
		if environment != "" {
			settings.environment = environment
		}

		return nil
	}
}

func WithServiceName(name string) ProviderOption {
	return func(settings *settings) error {
		if name != "" {
			settings.name = name
		}

		return nil
	}
}

func WithVersion(version string) ProviderOption {
	return func(settings *settings) error {
		if version != "" {
			settings.version = version
		}

		return nil
	}
}

func WithLogger(lgr logger) ProviderOption {
	return func(settings *settings) error {
		if lgr != nil {
			settings.lgr = &errorLogger{lgr: lgr}
		}

		return nil
	}
}

func WithSampleRatio(ratio float64) ProviderOption {
	return func(settings *settings) error {
		if ratio > 0 && ratio <= 1 {
			settings.sampleRatio = ratio
		}

		return nil
	}
}

func WithReceiverEndpoint(rcvr string) ProviderOption {
	return func(settings *settings) error {
		if rcvr != "" {
			settings.receiver = rcvr
		}

		return nil
	}
}

func WithExporterTimeout(timeout time.Duration) ProviderOption {
	return func(settings *settings) error {
		if timeout != 0 {
			settings.exporterTimeout = timeout
		}

		return nil
	}
}

func WithAttributesInjector(injector AttributeInjector) ProviderOption {
	return func(settings *settings) error {
		if injector == nil {
			return nil
		}

		if settings.traceInjectors == nil {
			settings.traceInjectors = []AttributeInjector{injector}
			return nil
		}

		settings.traceInjectors = append(settings.traceInjectors, injector)
		return nil
	}
}

func WithTextMapPropagtor(propagtor propagation.TextMapPropagator) ProviderOption {
	return func(settings *settings) error {
		if propagtor == nil {
			return nil
		}

		if settings.propagators == nil {
			settings.propagators = []propagation.TextMapPropagator{propagtor}
			return nil
		}

		settings.propagators = append(settings.propagators, propagtor)

		return nil
	}
}

func WithB3Propagtor() ProviderOption {
	return func(settings *settings) error {
		propagtor := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))

		if settings.propagators == nil {
			settings.propagators = []propagation.TextMapPropagator{propagtor}
			return nil
		}

		settings.propagators = append(settings.propagators, propagtor)

		return nil
	}
}
