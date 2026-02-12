// Adapter converts domain-level handlers into http.Handler.
// This allows business logic to remain framework-agnostic.

package api

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"golang.org/x/exp/slog"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) *APIError

// Adapter wraps a domain handler and converts returned APIError
// into a standardized HTTP response.
func Adapter(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if rec := recover(); rec != nil {
				reqID := GetRequestID(r.Context())

				slog.Error("panic recovered",
					"request_id", reqID,
					"method", r.Method,
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
					"panic", fmt.Sprintf("%v", rec),
					"stack", string(debug.Stack()),
				)

				WriteError(w, r, &APIError{
					Status:  http.StatusInternalServerError,
					Code:    ErrInternal,
					Message: "Internal server error",
				})
			}
		}()

		if err := h(w, r); err != nil {
			reqID := GetRequestID(r.Context())

			slog.Error("handler error",
				"request_id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", err.Status,
				"code", err.Code,
				"message", err.Message,
			)

			WriteError(w, r, err)
		}
	}
}
