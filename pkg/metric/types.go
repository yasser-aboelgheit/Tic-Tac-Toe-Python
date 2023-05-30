package metric

import "time"

// Counter represents counter metric.
type Counter struct {
	Name string
}

// Count tracks how many times this counter happened per second.
func (c *Counter) Count(value int64, tags []string, rate float64) {
	Count(c.Name, value, tags, rate)
}

// Incr is just Count +1.
func (c *Counter) Incr(tags []string, rate float64) {
	Incr(c.Name, tags, rate)
}

// Decr is just Count -1.
func (c *Counter) Decr(tags []string, rate float64) {
	Decr(c.Name, tags, rate)
}

// NamedGauge represents gauge with a name.
type NamedGauge struct {
	Name string
}

// Set sets gauge value.
func (ng *NamedGauge) Set(value float64, tags []string, rate float64) {
	Gauge(ng.Name, value, tags, rate)
}

// NamedTiming is a named timing.
type NamedTiming struct {
	Name string
}

// Timing reports duration of operation.
func (t *NamedTiming) Timing(dur time.Duration, tags []string, rate float64) {
	Timing(t.Name, dur, tags, rate)
}

// NamedDistribution is a named distribution.
type NamedDistribution struct {
	Name string
}

// Distribution tracks the statistical distribution of a set of values across your infrastructure.
func (h *NamedDistribution) Distribution(value float64, tags []string, rate float64) {
	Histogram(h.Name, value, tags, rate)
}
