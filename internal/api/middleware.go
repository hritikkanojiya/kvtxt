package api

import (
	"io"
	"log/slog"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(start).Milliseconds(),
			"remote_addr", r.RemoteAddr,
		)
	})
}

func MaxBodySize(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}

func AllowMethods(methods ...string) func(HandlerFunc) HandlerFunc {
	allowed := make(map[string]struct{}, len(methods))
	for _, m := range methods {
		allowed[m] = struct{}{}
	}

	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) *APIError {

			if _, ok := allowed[r.Method]; !ok {

				if r.Body != nil {
					io.Copy(io.Discard, r.Body)
					r.Body.Close()
				}

				return &APIError{
					Status:  http.StatusMethodNotAllowed,
					Code:    ErrMethodNotAllowed,
					Message: "Method not allowed",
				}
			}

			return next(w, r)
		}
	}
}
