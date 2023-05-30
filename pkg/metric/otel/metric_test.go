package otel

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func TestCreateMeterProvider(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		_, err := CreateMeterProvider(context.Background(), "localhost:4317")
		if err != nil {
			t.Fatalf("could not create provider: %v", err)
		}
	})

	t.Run("error option", func(t *testing.T) {
		t.Parallel()

		_, err := CreateMeterProvider(context.Background(), "localhost:4317", optionFunc(func(c *config) error {
			return fmt.Errorf("test err")
		}))
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})

	t.Run("empty receiver endpoint", func(t *testing.T) {
		t.Parallel()

		_, err := CreateMeterProvider(context.Background(), "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})

	t.Run("invalid receiver endpoint", func(t *testing.T) {
		t.Parallel()

		_, err := CreateMeterProvider(context.Background(), "!@#$%")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		t.Log(err)
	})
}

func Test_globalAttributes(t *testing.T) {
	t.Parallel()

	want := []attribute.KeyValue{
		semconv.DeploymentEnvironment("test-env"),
		semconv.ServiceName("test-svc"),
		semconv.ServiceVersion("0.0.1-test"),
	}

	got := globalAttributes(config{
		environment: "test-env",
		serviceName: "test-svc",
		version:     "0.0.1-test",
	})

	if !reflect.DeepEqual(want, got) {
		t.Error("attributes mismatch")
		t.Errorf("want: %+v", want)
		t.Errorf("got:  %+v", got)
	}
}

func Test_createExporter(t *testing.T) {
	t.Parallel()

	type args struct {
		useCumulativeTemporality bool
	}

	tests := []struct {
		name                string
		args                args
		wantHistTemporality metricdata.Temporality
	}{
		{
			name:                "cumulative temporality",
			args:                args{useCumulativeTemporality: true},
			wantHistTemporality: metricdata.CumulativeTemporality,
		},
		{
			name:                "delta temporality",
			args:                args{useCumulativeTemporality: false},
			wantHistTemporality: metricdata.DeltaTemporality,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exporter, err := createExporter(context.Background(), "localhost:4317", tt.args.useCumulativeTemporality)
			if err != nil {
				t.Fatalf("could not create exporter: %v", err)
			}

			gotHistTemporality := exporter.Temporality(metric.InstrumentKindHistogram)
			if tt.wantHistTemporality != gotHistTemporality {
				t.Fatalf("temporality mismatch: want %d; got %d", tt.wantHistTemporality, gotHistTemporality)
			}
		})
	}
}
