package otel

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

const (
	// exportTimeout defines how long metric exporter will attempt to export a batch of metrics.
	exportTimeout = 5 * time.Second

	// defaultExportInterval is the default interval of exporting metrics.
	defaultExportInterval = 15 * time.Second
)

var defaultConfig = config{
	exportInterval: defaultExportInterval,
}

// Closer closes metrics.
type Closer func(context.Context) error

func noopclose(context.Context) error { return nil }

// Run runs metrics without returning provider to caller.
func Run(ctx context.Context, receiverEndpoint string, opts ...Option) Closer {
	if receiverEndpoint == "" {
		return noopclose
	}

	meterProvider, err := CreateMeterProvider(
		ctx,
		receiverEndpoint,
		opts...,
	)
	if err != nil {
		log.Printf("could not create meter provider: %v\n", err)

		return noopclose
	}

	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown
}

// CreateMeterProvider constructs provider of meters.
func CreateMeterProvider(
	ctx context.Context,
	receiverEndpoint string,
	opts ...Option,
) (*metric.MeterProvider, error) {
	cfg := defaultConfig
	for _, opt := range opts {
		err := opt.apply(&cfg)
		if err != nil {
			return nil, fmt.Errorf("cannot apply an option: %w", err)
		}
	}

	// set up exporter
	if receiverEndpoint == "" {
		return nil, fmt.Errorf("metrics receiver endpoint not specified")
	}

	exporter, err := createExporter(ctx, receiverEndpoint, cfg.useCumulativeTemporality)
	if err != nil {
		return nil, err
	}
	reader := metric.NewPeriodicReader(exporter, metric.WithInterval(cfg.exportInterval),
		metric.WithTimeout(exportTimeout))

	// set up resource
	res, err := resource.New(ctx, resource.WithHost(), resource.WithAttributes(globalAttributes(cfg)...))
	if err != nil {
		return nil, fmt.Errorf("failed to create meter provider resource: %w", err)
	}

	// create provider
	provider := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))

	return provider, nil
}

func globalAttributes(cfg config) []attribute.KeyValue {
	var attrs []attribute.KeyValue

	if cfg.environment != "" {
		attrs = append(attrs, semconv.DeploymentEnvironment(cfg.environment))
	}

	if cfg.serviceName != "" {
		attrs = append(attrs, semconv.ServiceName(cfg.serviceName))
	}

	if cfg.version != "" {
		attrs = append(attrs, semconv.ServiceVersion(cfg.version))
	}

	attrs = append(attrs, cfg.attrs...)

	return attrs
}

func createExporter( //nolint:ireturn
	ctx context.Context,
	receiverEndpoint string,
	useCumulativeTemporality bool,
) (metric.Exporter, error) {
	temporality := deltaPreferredTemporalitySelector
	if useCumulativeTemporality {
		temporality = alwaysCumulativeTemporalitySelector
	}

	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithTemporalitySelector(temporality),
		otlpmetricgrpc.WithEndpoint(receiverEndpoint),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithTimeout(exportTimeout))
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}
	return exporter, nil
}

// alwaysCumulativeTemporalitySelector returns cumulative temporality for all
// instruments. Cumulative temporality is used by Prometheus.
func alwaysCumulativeTemporalitySelector(metric.InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
}

// deltaPreferredTemporalitySelector indicates delta aggregation temporality
// preference. Delta temporality is used by NewRelic, Datadog, etc. See
// https://docs.datadoghq.com/opentelemetry/guide/otlp_delta_temporality/?code-lang=go
func deltaPreferredTemporalitySelector(kind metric.InstrumentKind) metricdata.Temporality {
	switch kind {
	case metric.InstrumentKindCounter,
		metric.InstrumentKindHistogram,
		metric.InstrumentKindObservableGauge,
		metric.InstrumentKindObservableCounter:
		return metricdata.DeltaTemporality
	case metric.InstrumentKindUpDownCounter,
		metric.InstrumentKindObservableUpDownCounter:
		return metricdata.CumulativeTemporality
	default:
		return metricdata.DeltaTemporality
	}
}
