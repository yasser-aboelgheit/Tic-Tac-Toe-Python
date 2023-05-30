package tracer

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type TraceProvider interface {
	Tracer(name string) Tracer
	Stop()
}

type tracerProvider struct {
	internalProvider *tracesdk.TracerProvider
	propagators      propagation.TextMapPropagator
	logger           logger
	resourceTimeout  time.Duration
}

func NewTraceProvider(ctx context.Context, opts ...ProviderOption) (*tracerProvider, error) {
	stngs := newDefaultSettings()
	for _, opt := range opts {
		err := opt(&stngs)
		if err != nil {
			return nil, fmt.Errorf("option error: %w", err)
		}
	}

	exporter, err := createExporter(ctx, stngs.receiver, stngs.exporterTimeout)
	if err != nil {
		return nil, fmt.Errorf("unable to create exporter: %w", err)
	}

	res, err := resource.New(
		ctx,
		resource.WithTelemetrySDK(),
		resource.WithContainer(),
		resource.WithContainerID(),
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.ServiceName(stngs.name),
			semconv.ServiceVersion(stngs.version),
			semconv.DeploymentEnvironment(stngs.environment),
		),
		resource.WithAttributes(stngs.attributes...),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create trace provider resource: %w", err)
	}

	rawTraceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(stngs.sampleRatio))),
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
	)

	propagator := propagation.NewCompositeTextMapPropagator(stngs.propagators...)

	otel.SetTracerProvider(rawTraceProvider)
	otel.SetTextMapPropagator(propagator)
	otel.SetErrorHandler(stngs.lgr)

	traceProvider := &tracerProvider{
		internalProvider: rawTraceProvider,
		propagators:      propagator,
		resourceTimeout:  stngs.resourceTimeout,
		logger:           stngs.lgr.lgr,
	}

	return traceProvider, nil
}

func (tp *tracerProvider) Tracer(scopeName string) Tracer {
	t := tp.provider().Tracer(scopeName)

	return NewTracer(t, tp)
}

func (tp *tracerProvider) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), tp.resourceTimeout)
	defer cancel()

	err := tp.internalProvider.Shutdown(ctx)
	if err != nil {
		tp.logger.Errorw(ctx, err, "error shutting down tracer", nil)
	}
}

func (tp *tracerProvider) provider() trace.TracerProvider {
	return tp.internalProvider
}

func (tp *tracerProvider) propagator() propagation.TextMapPropagator {
	return tp.propagators
}
