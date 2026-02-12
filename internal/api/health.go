package api

import "net/http"

func Health() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *APIError {
		WriteJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
		return nil
	}
}
