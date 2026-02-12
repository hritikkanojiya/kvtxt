// Package api contains HTTP layer logic including routing,
// middleware, request/response models, and error handling.
//
// This file defines all HTTP route registrations and
// composes middleware chain.

package api

import (
	"net/http"
)

// RegisterRoute registers all HTTP endpoints.
// Each route should delegate to a thin handler that calls service/storage layer.
// Avoid embedding business logic inside handlers.
func RegisterRoute(
	mux *http.ServeMux,
	path string,
	method string,
	h HandlerFunc,
) {
	handler := Adapter(
		AllowHttpMethods(method)(h),
	)

	mux.Handle(path, handler)
}
