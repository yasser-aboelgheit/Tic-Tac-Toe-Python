package otel

import (
	"reflect"
	"testing"
	"time"
)

func TestApplyOptions(t *testing.T) {
	t.Parallel()

	opts := []Option{
		WithExportInterval(time.Second * 42),
		WithServiceNameAttribute("test-svc"),
		WithEnvironmentAttribute("test-env"),
		WithVersionAttribute("0.0.1-test"),
		WithCumulativeTemporality(),
	}

	got := config{}
	for _, opt := range opts {
		err := opt.apply(&got)
		if err != nil {
			t.Fatalf("could not apply option: %v", err)
		}
	}

	want := config{
		exportInterval:           time.Second * 42,
		serviceName:              "test-svc",
		environment:              "test-env",
		version:                  "0.0.1-test",
		useCumulativeTemporality: true,
	}

	if !reflect.DeepEqual(want, got) {
		t.Error("config mismatch")
		t.Errorf("want: %+v", want)
		t.Errorf("got:  %+v", got)
	}
}
