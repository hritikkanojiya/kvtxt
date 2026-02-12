package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/hritikkanojiya/kvtxt/internal/constant"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}

		w.Header().Set("X-Request-ID", reqID)

		ctx := context.WithValue(r.Context(), constant.RequestIdKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(constant.RequestIdKey).(string); ok {
		return v
	}
	return ""
}
