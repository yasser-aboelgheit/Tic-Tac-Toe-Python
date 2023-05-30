package logger

import (
	"context"

)

// logcontextKeys is a namespace for logger related attributes within context.
type logContextKey struct{}

// ContextualAttributes is a container for logger related attributes.
type contextAttributes struct {
	attrs map[string]interface{}
}

func readContextAttributes(ctx context.Context) (map[string]any, bool) {
	if ctx == nil {
		return map[string]any{}, false
	}

	ctxReader := ctx.Value(logContextKey{})
	if attrs, ok := ctxReader.(contextAttributes); ok {
		return attrs.attrs, true
	}

	return map[string]any{}, false
}

// ContextWithAttributes responsible to add loggerable attributes to context.
func ContextWithAttributes(ctx context.Context, attrs map[string]any) context.Context {
	if attrs == nil {
		return ctx
	}

	if oldAttrs, ok := readContextAttributes(ctx); ok {
		attrs = mergeAttr(oldAttrs, attrs)
	}

	return context.WithValue(ctx, logContextKey{}, createContextAttributes(attrs))
}

// createContextAttributes make a new copy of a dict.
func createContextAttributes(original map[string]any) contextAttributes {
	dest := make(map[string]interface{}, len(original))
	for i, v := range original {
		dest[i] = v
	}
	return contextAttributes{attrs: dest}
}

// mergeAttr returns new contextAttributes with both keys.
// NOTE: attrs2 will override any same key in attrs1.
func mergeAttr(attrs1, attrs2 map[string]any) map[string]any {
	result := make(map[string]any, len(attrs1)+len(attrs2))

	for i, v := range attrs1 {
		result[i] = v
	}

	for i, v := range attrs2 {
		result[i] = v
	}

	return result
}
