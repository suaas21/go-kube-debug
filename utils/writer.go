package utils

import "net/http"

type ResponseWriterInterceptor struct {
	http.ResponseWriter
	StatusCode int
}

func (w *ResponseWriterInterceptor) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
