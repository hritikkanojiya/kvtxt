package api

import (
	"io"
	"net/http"
)

func AllowHttpMethods(methods ...string) func(HandlerFunc) HandlerFunc {
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
