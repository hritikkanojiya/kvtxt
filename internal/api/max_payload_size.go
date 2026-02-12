// MaxBodySize limits incoming request payload size to prevent
// memory exhaustion and abuse.
//
// If the payload exceeds configured size, a 413 error is returned.

package api

import "net/http"

func MaxPayloadSize(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Wrap request body with MaxBytesReader to enforce hard limit
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}
