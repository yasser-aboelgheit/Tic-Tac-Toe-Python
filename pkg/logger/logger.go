// Package logger provide a logger definition that can log with all different levels
// "DEBUG", "INFO", "WARN" and "ERROR" with addition to add all the kv needed to be attached,
// all keys attached to the context will be logged as well, unless

package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger interface {
	Infow(ctx context.Context, msg string, attrs map[string]interface{})
	Debugw(ctx context.Context, msg string, attrs map[string]interface{})
	Warnw(ctx context.Context, msg string, attrs map[string]interface{})
	Errorw(ctx context.Context, err error, msg string, attrs map[string]interface{})
	Info(ctx context.Context, msg string)
	Debug(ctx context.Context, msg string)
	Warn(ctx context.Context, msg string)
	Error(ctx context.Context, err error, msg string)
	Println(...interface{})
}

// Level logging level.
type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"

	serviceKey           string = "service"
	versionKey           string = "version"
	environmentKey       string = "env"
	minLogLevelKey       string = "minimum_log_level"
	externalCallerOffset int    = 2
)

// Log wraps whatever logging service we want to use.
type Log struct {
	root                 zerolog.Logger
	externalCallerOffset int
}

func init() { //nolint:gochecknoinits
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

func NewLogger(cfg *Config, opts ...Option) Log {
	cfg.apply(opts...)

	lgr := log.Logger.With().
		Str(serviceKey, cfg.serviceName).
		Str(versionKey, cfg.version).
		Str(environmentKey, cfg.environemnt).
		Str(minLogLevelKey, string(cfg.Level)).
		Logger().
		Level(parseLogLevel(string(cfg.Level)))

	localLgr := Log{root: lgr, externalCallerOffset: externalCallerOffset}

	if cfg.PrettyPrint {
		return localLgr.WithPrettyOutput()
	}

	return localLgr
}

// Infow is like log.Info level but with attributes to be added.
func (l *Log) Infow(
	ctx context.Context,
	msg string,
	attrs map[string]interface{},
) {
	l.createEvent(ctx, zerolog.InfoLevel, attrs).Msg(msg)
}

// Debugw is like log.Debug level but with attributes to be added.
func (l *Log) Debugw(
	ctx context.Context,
	msg string,
	attrs map[string]interface{},
) {
	l.createEvent(ctx, zerolog.DebugLevel, attrs).Msg(msg)
}

// Warnw is like log.Warn level but with attributes to be added.
func (l *Log) Warnw(
	ctx context.Context,
	msg string,
	attrs map[string]interface{},
) {
	l.createEvent(ctx, zerolog.WarnLevel, attrs).Msg(msg)
}

// Errorw is like log.Error level but with attributes to be added.
func (l *Log) Errorw(
	ctx context.Context,
	err error,
	msg string,
	attrs map[string]interface{},
) {
	l.createEvent(ctx, zerolog.ErrorLevel, attrs).Err(err).Msg(msg)
}

// Info is like log.Info level but all previous attributes.
func (l *Log) Info(ctx context.Context, msg string) {
	l.createEvent(ctx, zerolog.InfoLevel, nil).Msg(msg)
}

// Debug is like log.Debug level but all previous attributes.
func (l *Log) Debug(ctx context.Context, msg string) {
	l.createEvent(ctx, zerolog.DebugLevel, nil).Msg(msg)
}

// Warn is like log.Warn level but all previous attributes.
func (l *Log) Warn(ctx context.Context, msg string) {
	l.createEvent(ctx, zerolog.WarnLevel, nil).Msg(msg)
}

// Error is like log.Error level but all previous attributes.
func (l *Log) Error(ctx context.Context, err error, msg string) {
	l.createEvent(ctx, zerolog.ErrorLevel, nil).Err(err).Msg(msg)
}

// createEvent to create Event whatever the logging type is.
func (l *Log) createEvent(
	ctx context.Context,
	level zerolog.Level,
	attrs map[string]interface{},
) *zerolog.Event {
	if l == nil {
		panic("no logger provided")
	}
	event := enrichEventWithCtx(ctx, l.root.WithLevel(level))
	if attrs != nil {
		event.Fields(attrs)
	}

	return event.Caller(l.externalCallerOffset)
}

func enrichEventWithCtx(
	ctx context.Context,
	event *zerolog.Event,
) *zerolog.Event {
	if ctxAttrs, ok := readContextAttributes(ctx); ok {
		event.Fields(ctxAttrs)
	}

	return event
}

// WithAttributes returns a new logger adding attributes.
func (l Log) WithAttributes(attrs map[string]interface{}) Log {
	l.root = l.root.With().Fields(attrs).Logger()

	return l
}

func (l Log) WithPrettyOutput() Log {
	return l.WithOutput(
		zerolog.ConsoleWriter{
			Out: os.Stderr,
			PartsExclude: []string{
				zerolog.CallerFieldName,
			},
			FieldsExclude: []string{
				minLogLevelKey,
				environmentKey,
				versionKey,
			},
		},
	)
}

// WithOutput returns a new logger with writer provided.
func (l Log) WithOutput(w io.Writer) Log {
	l.root = l.root.Output(w)
	return l
}

func (l Log) Println(values ...interface{}) {
	l.Error(context.Background(), fmt.Errorf("%v", values), "panic occurred")
}

// parseLogLevel translate log level from string to zerolog.level.
func parseLogLevel(level string) zerolog.Level {
	zlLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		zlLevel = zerolog.InfoLevel
	}

	return zlLevel
}
