package middleware

import (
	"context"
	"net/http"
	"time"
)

// MakeTimeoutMiddleware adds request maxRequestDuration to request context.
func MakeTimeoutMiddleware(maxRequestDuration time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxWithTimeout, cancelFunc := context.WithTimeout(r.Context(), maxRequestDuration)
			defer cancelFunc()

			newRequest := r.WithContext(ctxWithTimeout)
			next.ServeHTTP(w, newRequest)
		})
	}
}
