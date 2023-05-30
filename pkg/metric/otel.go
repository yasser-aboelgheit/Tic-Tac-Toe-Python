package metric

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// SetMeterProvider sets global meter provider.
func SetMeterProvider(meterprovider metric.MeterProvider) {
	otel.SetMeterProvider(meterprovider)
}

// GetMeterProvider returns global meter provider.
func GetMeterProvider() metric.MeterProvider {
	return otel.GetMeterProvider()
}

// Meter creates meter from global meter provider.
func Meter(name string) metric.Meter {
	return otel.Meter(name)
}
