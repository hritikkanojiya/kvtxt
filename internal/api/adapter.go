package api

import "net/http"

type HandlerFunc func(w http.ResponseWriter, r *http.Request) *APIError

func Adapt(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if rec := recover(); rec != nil {
				WriteError(w, r, &APIError{
					Status:  http.StatusInternalServerError,
					Code:    ErrInternal,
					Message: "Internal server error",
				})
			}
		}()

		if err := h(w, r); err != nil {
			WriteError(w, r, err)
		}
	}
}
