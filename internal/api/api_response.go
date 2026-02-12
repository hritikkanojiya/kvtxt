// APIResponse defines the standard success response format.
// All successful responses must use this structure to maintain API consistency.

package api

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error struct {
		Code    ErrorCode `json:"code"`
		Message string    `json:"message"`
	} `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, r *http.Request, err *APIError) {
	reqID := GetRequestID(r.Context())

	resp := errorResponse{}
	resp.Error.Code = err.Code
	resp.Error.Message = err.Message

	if reqID != "" {
		w.Header().Set("X-Request-ID", reqID)
	}

	WriteJSON(w, err.Status, resp)
}
