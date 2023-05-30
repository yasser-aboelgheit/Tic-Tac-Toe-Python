package tracer

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

const (
	resourceNameKey string = "resource.name"
	spanIDKey       string = "otel.span_id"
	traceIDKey      string = "otel.trace_id"
)

func createExporter(
	ctx context.Context,
	receiver string,
	timeout time.Duration,
) (*otlptrace.Exporter, error) {
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(receiver),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(timeout),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("create trace exporter error: %w", err)
	}

	return exporter, nil
}

func mapToAttributes(attributesMap map[string]any) []attribute.KeyValue {
	var attributes []attribute.KeyValue

	for key, value := range attributesMap {
		var attrVal attribute.Value
		switch val := value.(type) {
		case string:
			attrVal = attribute.StringValue(val)
		case int:
			attrVal = attribute.IntValue(val)
		case int32:
		case int64:
			attrVal = attribute.Int64Value(val)
		case float32:
		case float64:
			attrVal = attribute.Float64Value(val)
		case bool:
			attrVal = attribute.BoolValue(val)
		default:
			attrVal = attribute.StringValue(fmt.Sprint(val))
		}

		attributes = append(attributes, attribute.KeyValue{
			Key:   attribute.Key(key),
			Value: attrVal,
		})
	}

	return attributes
}
