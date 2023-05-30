package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	metricsPath = flag.String("metrics", "", "metrics configuration file")
	packageName = flag.String("package", "", "generated file package name")
	outputPath  = flag.String("output", "", "output file name; default <config_dir>/metrics_gen.go")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of metricgen:\n")
	fmt.Fprintf(os.Stderr, "\t metricgen -config config.yaml -package metrics\n")
	fmt.Fprintf(os.Stderr, "\t metricgen -config config.yaml -package metrics -output ../metrics/metrics.go\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("metricgen: ")

	flag.Usage = Usage
	flag.Parse()

	if *metricsPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	if *packageName == "" {
		flag.Usage()
		os.Exit(2)
	}

	output, err := getOutput(*outputPath, *metricsPath)
	if err != nil {
		log.Fatalf("could not get output path: %v", err)
	}

	metrics, err := LoadMetrics(*metricsPath)
	if err != nil {
		log.Fatalf("could not load metrics: %v", err)
	}

	if len(metrics) == 0 {
		log.Printf("config %q has no metrics", *metricsPath)
		return
	}

	g := generator{
		Args:    strings.Join(os.Args[1:], " "),
		Package: *packageName,
	}

	if err = g.processMetrics(metrics); err != nil {
		log.Fatalf("could not process metrics: %v", err)
	}

	source, err := g.generate()
	if err != nil {
		log.Fatalf("could not generate metrics: %v", err)
	}

	if err := write(output, source); err != nil {
		log.Fatalf("could not write generated metrics: %v", err)
	}
}

// getOutput returns output filename.
func getOutput(path, fallback string) (string, error) {
	const defaultname = "metrics_gen.go"

	if path == "" {
		output := filepath.Dir(fallback)
		output = filepath.Join(output, defaultname)

		return filepath.Abs(output)
	}

	path = filepath.Clean(path)

	var filename string
	{
		filename = filepath.Base(path)
		if filename == "." {
			filename = defaultname
		}

		ext := filepath.Ext(filename)
		if ext != ".go" {
			filename = defaultname
		}
	}

	dir := filepath.Dir(path)

	return filepath.Abs(filepath.Join(dir, filename))
}

func write(output string, source []byte) error {
	dst, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("could not open output: %w", err)
	}

	if _, err := io.Copy(dst, bytes.NewReader(source)); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}

	return nil
}
