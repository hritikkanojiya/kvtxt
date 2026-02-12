package api

import (
	"net/http"

	"github.com/hritikkanojiya/kvtxt/internal/storage"
)

func Readiness(store *storage.Storage) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *APIError {

		if err := store.Ping(); err != nil {
			return &APIError{
				Status:  http.StatusServiceUnavailable,
				Code:    ErrInternal,
				Message: "Storage not ready",
			}
		}

		WriteJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})

		return nil
	}
}
