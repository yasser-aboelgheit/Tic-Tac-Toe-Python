// Package metrics allows to send metrics values.
package metric

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

var stats *statsd.Client

// DefaultAddr returns default dogstatsd addr.
func DefaultAddr() string { return os.Getenv("METRICS_ENDPOINT") }

// Options control metrics options.
type Options struct {
	// Tags control metric tags.
	Tags []string
}

// Option overrides metrics options.
type Option func(*Options) error

// WithTags adds tags to all reported metrics.
func WithTags(tags []string) Option {
	return func(opts *Options) error {
		opts.Tags = tags
		return nil
	}
}

// Run runs metric pushing.
func Run(addr string, opts ...Option) (err error) {
	sopts, err := resolveOptions(opts...)
	if err != nil {
		return fmt.Errorf("could not resolve options: %w", err)
	}

	stats, err = statsd.New(addr, sopts...)
	if err != nil {
		return fmt.Errorf("could not initialize stats: %w", err)
	}

	return nil
}

func resolveOptions(opts ...Option) ([]statsd.Option, error) {
	if len(opts) == 0 {
		return nil, nil
	}

	var options Options
	for _, opt := range opts {
		if opt == nil {
			return nil, errors.New("expected Option but got nil")
		}

		if err := opt(&options); err != nil {
			return nil, err
		}
	}

	result := make([]statsd.Option, 0, len(opts))

	if len(options.Tags) > 0 {
		result = append(result, statsd.WithTags(options.Tags))
	}

	return result, nil
}

// Stop stops metric pushing.
func Stop() {
	stats.Close() //nolint:errcheck
}

// Gauge measures the value of a metric at a particular time.
func Gauge(name string, value float64, tags []string, rate float64) {
	stats.Gauge(name, value, tags, rate) //nolint:errcheck
}

// Count tracks how many times something happened per second.
func Count(name string, value int64, tags []string, rate float64) {
	stats.Count(name, value, tags, rate) //nolint:errcheck
}

// Histogram tracks the statistical distribution of a set of values on each host.
func Histogram(name string, value float64, tags []string, rate float64) {
	stats.Histogram(name, value, tags, rate) //nolint:errcheck
}

// Distribution tracks the statistical distribution of a set of values across your infrastructure.
func Distribution(name string, value float64, tags []string, rate float64) {
	stats.Distribution(name, value, tags, rate) //nolint:errcheck
}

// Decr is just Count of -1
func Decr(name string, tags []string, rate float64) {
	stats.Decr(name, tags, rate) //nolint:errcheck
}

// Incr is just Count of 1
func Incr(name string, tags []string, rate float64) {
	stats.Incr(name, tags, rate) //nolint:errcheck
}

// Set counts the number of unique elements in a group.
func Set(name string, value string, tags []string, rate float64) {
	stats.Set(name, value, tags, rate) //nolint:errcheck
}

// Timing sends timing information, it is an alias for TimeInMilliseconds
func Timing(name string, value time.Duration, tags []string, rate float64) {
	stats.Timing(name, value, tags, rate) //nolint:errcheck
}

// TimeInMilliseconds sends timing information in milliseconds.
// It is flushed by statsd with percentiles, mean and other info
// (https://github.com/etsy/statsd/blob/master/docs/metric_types.md#timing)
func TimeInMilliseconds(name string, value float64, tags []string, rate float64) {
	stats.TimeInMilliseconds(name, value, tags, rate) //nolint:errcheck
}

// Flush forces a flush of all the queued dogstatsd payloads.
func Flush() error {
	if err := stats.Flush(); err != nil {
		return fmt.Errorf("could not flush stats: %w", err)
	}

	return nil
}
