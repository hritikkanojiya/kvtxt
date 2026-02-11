package api

import (
	"net/http"
)

func RegisterRoute(
	mux *http.ServeMux,
	path string,
	method string,
	h HandlerFunc,
) {
	handler := Adapt(
		AllowMethods(method)(h),
	)

	mux.Handle(path, handler)
}
