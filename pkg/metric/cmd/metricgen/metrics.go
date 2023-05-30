package main

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

// Metrics is a collection of metrics.
type Metrics struct {
	Metrics []Metric `yaml:"metrics"`
}

// Metric represents custom metric.
type Metric struct {
	// Name is a metric name.
	Name string `yaml:"name"`
	// Type is a metric type.
	Type MetricType `yaml:"type"`
	// ShortName is a name used in generated code.
	ShortName string `yaml:"short_name"`
}

// MetricType is a custom metric type.
type MetricType string

const (
	// TypeCounter stands for counter metric type.
	TypeCounter MetricType = "counter"
	// TypeGauge stands for gague metric type.
	TypeGauge MetricType = "gauge"
	// TypeTiming stands for timings metric type.
	TypeTiming MetricType = "timing"
	// TypeDistribution stands for distribution metric type.
	TypeDistribution MetricType = "distribution"
)

// LoadMetrics reads metrics configuration file and returns metrics.
func LoadMetrics(path string) ([]Metric, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	var metrics Metrics
	if err := yaml.NewDecoder(f).Decode(&metrics); err != nil {
		return nil, fmt.Errorf("could not read metrics: %w", err)
	}

	return metrics.Metrics, nil
}
