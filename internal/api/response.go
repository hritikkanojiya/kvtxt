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

func WriteError(w http.ResponseWriter, err *APIError) {
	resp := errorResponse{}
	resp.Error.Code = err.Code
	resp.Error.Message = err.Message

	WriteJSON(w, err.Status, resp)
}
