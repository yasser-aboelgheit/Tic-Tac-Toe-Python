package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"text/template"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	//go:embed template.tmpl
	metricsraw  string
	metricstmpl = template.Must(template.New("metrics").Parse(metricsraw))

	caser = cases.Title(language.English)
)

type generator struct {
	Package       string
	Args          string
	Counters      []simpleMetric
	Gauges        []simpleMetric
	Timings       []simpleMetric
	Distributions []simpleMetric
}

type simpleMetric struct {
	Name       string
	MetricName string
}

func (g *generator) generate() ([]byte, error) {
	var buf bytes.Buffer
	if err := metricstmpl.Execute(&buf, g); err != nil {
		return nil, fmt.Errorf("could not execute template: %w", err)
	}

	source, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("could not format generated code: %w", err)
	}

	return source, nil
}

func (g *generator) processMetrics(metrics []Metric) error {
	for _, m := range metrics {
		if err := g.processMetric(m); err != nil {
			return fmt.Errorf("could not process metric: %w", err)
		}
	}

	return nil
}

func (g *generator) processMetric(metric Metric) error {
	switch metric.Type {
	case TypeCounter:
		g.Counters = append(g.Counters, simpleMetric{
			Name:       getname(metric),
			MetricName: metric.Name,
		})
		return nil
	case TypeGauge:
		g.Gauges = append(g.Gauges, simpleMetric{
			Name:       getname(metric),
			MetricName: metric.Name,
		})
		return nil
	case TypeTiming:
		g.Timings = append(g.Timings, simpleMetric{
			Name:       getname(metric),
			MetricName: metric.Name,
		})
		return nil

	case TypeDistribution:
		g.Distributions = append(g.Distributions, simpleMetric{
			Name:       getname(metric),
			MetricName: metric.Name,
		})
	}

	return fmt.Errorf("unkwnown metric type: %q", metric.Type)
}

func getname(metric Metric) string {
	if metric.ShortName != "" {
		return metric.ShortName
	}

	return normalize(metric.Name) + getTypeSuffix(metric.Type)
}

func getTypeSuffix(typ MetricType) string {
	suffix := string(typ)
	return caser.String(suffix)
}

func normalize(name string) string {
	var result string

	capitalize := true
	for _, c := range name {
		switch c {
		case ' ', '-', '_', '.':
			capitalize = true
			continue
		}

		if capitalize {
			c = unicode.ToUpper(c)
			capitalize = false
		}

		result = result + string(c)
	}

	return result
}
