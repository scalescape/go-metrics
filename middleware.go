package metrics

import (
	"log"
	"net/http"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.StatusCode = statusCode
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, StatusCode: 200}
}

type Labels struct {
	Service string
	Method  string
	Path    string
	Status  string
}

type Config struct {
	Address     string `split_words:"true" required:"true"`
	ServiceID   string `split_words:"true" default:"service"`
	labels      Labels
	enablePprof bool
}

type Observer struct {
	Config
	latencyMetric *LatencyMetric
	errorMetric   *ErrorMetric
}

func (o Observer) Middleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		mw := NewResponseWriter(w)
		next.ServeHTTP(mw, r)
		latency := time.Since(start)

		if err := o.latencyMetric.Capture(r, mw, latency); err != nil {
			log.Printf("E! error collecting latency metric: %v", err)
			return
		}

		if err := o.errorMetric.Capture(r, mw); err != nil {
			log.Printf("E! error collecting count metric: %v", err)
			return
		}
	}
	return http.HandlerFunc(f)
}

func (o Observer) Close() {
}
