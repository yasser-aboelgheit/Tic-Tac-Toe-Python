// Provide httpHandler middleware

package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type requestLogger interface {
	Infow(ctx context.Context, msg string, attrs map[string]interface{})
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *map[string]any
}

func newLogResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{ResponseWriter: w}
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *logResponseWriter) Write(data []byte) (int, error) {
	var body map[string]any
	json.Unmarshal(data, &body)
	// TODO: mask body response to remove hash sensitive data.
	w.body = &body
	return w.ResponseWriter.Write(data)
}

func MakeMuxLoggerMiddleware(lgr requestLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logResponseWriter := newLogResponseWriter(w)

			defer func(startTime time.Time) {
				// serve the route
				if logResponseWriter.statusCode == 0 {
					logResponseWriter.statusCode = 200
				}

				kv := map[string]interface{}{
					"transport":      "http",
					"status":         logResponseWriter.statusCode,
					"duration_human": time.Since(startTime).String(),
					"duration":       time.Since(startTime),
					"response": map[string]any{
						"headers": w.Header(),
						"body":    logResponseWriter.body,
					},
					"request": map[string]interface{}{
						"name":   r.URL.Path,
						"method": r.Method,
						"client": r.Header,
					},
				}

				lgr.Infow(
					r.Context(),
					"processed request",
					kv,
				)
			}(time.Now())

			next.ServeHTTP(logResponseWriter, r)
		})
	}
}
