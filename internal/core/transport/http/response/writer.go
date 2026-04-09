package core_http_response

import (
	"net/http"
)

var (
	StatusCodeNonInitialized = -1
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(r http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: r,
		statusCode:     StatusCodeNonInitialized,
	}
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.statusCode = statusCode
}

func (r *ResponseWriter) GetStatusCodeOrPanic() int {
	if r.statusCode == StatusCodeNonInitialized {
		panic("no status code set!")
	}

	return r.statusCode
}
