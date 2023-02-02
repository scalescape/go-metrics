package common

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.StatusCode = statusCode
}

func (rw *ResponseWriter) Status() int {
	return rw.StatusCode
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, StatusCode: 200}
}
