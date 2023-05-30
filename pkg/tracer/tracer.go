package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	StartSpan(ctx context.Context, name string, attr map[string]any) (context.Context, Span)
}

type tracer struct {
	internalTracer   trace.Tracer
	internalProvider TraceProvider
}

func NewTracer(trace trace.Tracer, provider TraceProvider) *tracer {
	return &tracer{
		internalTracer:   trace,
		internalProvider: provider,
	}
}

func (t *tracer) StartSpan(
	ctx context.Context,
	scopeName string,
	attr map[string]any,
) (context.Context, Span) {
	ctx, otelSpan := t.internalTracer.Start(
		ctx,
		scopeName,
		trace.WithAttributes(mapToAttributes(attr)...),
		trace.WithSpanKind(trace.SpanKindInternal),
	)

	return ctx, &span{span: otelSpan}
}

func (t *tracer) provider() TraceProvider {
	return t.internalProvider
}
