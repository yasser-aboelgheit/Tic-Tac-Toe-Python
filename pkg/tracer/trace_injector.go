package tracer

import (
	"context"
)

type AttributeInjector interface {
	Inject(ctx context.Context, traceAttr map[string]any) context.Context
}

type InjectAttributesFunc func(context.Context, map[string]any) context.Context

func (f InjectAttributesFunc) Inject(ctx context.Context, traceAttr map[string]any) context.Context {
	return f(ctx, traceAttr)
}
