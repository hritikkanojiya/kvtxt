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
	handler := Adapter(
		AllowHttpMethods(method)(h),
	)

	mux.Handle(path, handler)
}
