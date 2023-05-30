package tracer

import "context"

type logger interface {
	Infow(ctx context.Context, msg string, attrs map[string]interface{})
	Debugw(ctx context.Context, msg string, attrs map[string]interface{})
	Warnw(ctx context.Context, msg string, attrs map[string]interface{})
	Errorw(ctx context.Context, err error, msg string, attrs map[string]interface{})
}

type nooplogger struct{}

func (lgr *nooplogger) Infow(ctx context.Context, msg string, attrs map[string]interface{})  {}
func (lgr *nooplogger) Debugw(ctx context.Context, msg string, attrs map[string]interface{}) {}
func (lgr *nooplogger) Warnw(ctx context.Context, msg string, attrs map[string]interface{})  {}
func (lgr *nooplogger) Errorw(ctx context.Context, err error, msg string, attrs map[string]interface{}) {
}

type errorLogger struct {
	lgr logger
}

func (errlgr *errorLogger) Handle(err error) {
	errlgr.lgr.Errorw(context.Background(), err, "tracer error", nil)
}
