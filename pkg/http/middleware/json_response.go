package middleware

import (
	"net/http"
	"strings"
)

const (
	jsonHeader           string = "application/json"
	mimeContentTypeJson  string = "application/json; charset=utf-8"
	contentTypeHeaderKey string = "Content-Type"
)

func MakeJsonContentMiddleware(enforce bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				contentType := r.Header.Get(contentTypeHeaderKey)
				jsonRequested := strings.Contains(contentType, jsonHeader)
				w.Header().Add(contentTypeHeaderKey, mimeContentTypeJson)

				if enforce && !jsonRequested {
					w.WriteHeader(http.StatusUnsupportedMediaType)
					return
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}
