package middleware

import (
	"context"
	"net/http"
)

type logger interface {
	Infow(ctx context.Context, msg string, attrs map[string]interface{})
	WithAttributes(map[string]interface{})
}

func MakeContextEnrichMiddleware(
	contextAdder func(
		ctx context.Context,
		attrs map[string]interface{},
	) context.Context,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				newRequest := r.WithContext(
					contextAdder(
						r.Context(),
						GetRequestInfo(r),
					),
				)
				next.ServeHTTP(w, newRequest)
			},
		)
	}
}


// GetRequestInfo give all request parameters needed to add to log/trace.
func GetRequestInfo(r *http.Request) map[string]interface{} {
	const (
		sessionID = "session_id"
		sessionIDHeader = "X-SESSION-ID"

		deviceTokenHeader = "X-DEVICE-HEADER"
		deviceID = "device_id"

		userAgentHeader = "User-Agent"
		userAgentField = "user_agent"
	)
	userAgent := r.Header.Get(userAgentHeader)
	return map[string]interface{}{

		sessionID: r.Header.Get(sessionIDHeader),
		userAgentField: userAgent,
	}
}
