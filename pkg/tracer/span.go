package tracer

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Span interface {
	SetAttributes(attributes map[string]any)
	End(err error)
}

type span struct {
	span trace.Span
}

// SetAttributes sets attributes of the Span.
func (s span) SetAttributes(attributes map[string]any) {
	s.span.SetAttributes(mapToAttributes(attributes)...)
}

// End completes the Span.
func (s span) End(err error) {
	defer s.span.End()

	if err != nil {
		s.span.RecordError(err)
		s.span.SetStatus(codes.Error, err.Error())
	}
}

// SpanFromContext returns the current Span from ctx.
// If no Span is currently set in ctx an implementation of a Span that performs no operations is returned.
func SpanFromContext(ctx context.Context) Span {
	sp := trace.SpanFromContext(ctx)

	return span{span: sp}
}
