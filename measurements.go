package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var DefaultLabels = Labels{
	Service: "service",
	Method:  "method",
	Path:    "path",
	Status:  "status",
}

type LatencyMetric struct {
	prometheus.ObserverVec
}

type ErrorMetric struct {
	*prometheus.GaugeVec
}

func (l *LatencyMetric) labels(r *http.Request, mw *ResponseWriter) map[string]string {
	return map[string]string{
		DefaultLabels.Method: r.Method,
		DefaultLabels.Path:   r.URL.Path,
		//TODO: mask it as 4xx,5xx,...
		DefaultLabels.Status: fmt.Sprintf("%d", mw.StatusCode),
	}
}

func (l *LatencyMetric) Capture(r *http.Request, mw *ResponseWriter, latency time.Duration) error {
	metric, err := l.GetMetricWith(l.labels(r, mw))
	if err != nil {
		return err
	}
	metric.Observe(float64(latency.Milliseconds()))
	return nil
}

func (m *ErrorMetric) labels(r *http.Request, mw *ResponseWriter) map[string]string {
	return map[string]string{
		DefaultLabels.Method: r.Method,
		DefaultLabels.Path:   r.URL.Path,
		DefaultLabels.Status: fmt.Sprintf("%d", mw.StatusCode),
	}
}

func (m *ErrorMetric) Capture(r *http.Request, mw *ResponseWriter) error {
	metric, err := m.GetMetricWith(m.labels(r, mw))
	if err != nil {
		return err
	}
	metric.Inc()
	return nil
}

func NewLatencyMetric(cfg Config) (*LatencyMetric, error) {
	latencyMetric := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "http_latency_ms"},
		[]string{DefaultLabels.Service, DefaultLabels.Method, DefaultLabels.Path, DefaultLabels.Status},
	)
	if err := prometheus.Register(latencyMetric); err != nil {
		return nil, err
	}
	lm, err := latencyMetric.CurryWith(map[string]string{DefaultLabels.Service: cfg.ServiceID})
	if err != nil {
		return nil, err
	}
	return &LatencyMetric{lm}, nil
}

func NewErrorMetric(cfg Config) (*ErrorMetric, error) {
	errorMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "http_requests"},
		[]string{DefaultLabels.Service, DefaultLabels.Method, DefaultLabels.Path, DefaultLabels.Status})
	if err := prometheus.Register(errorMetric); err != nil {
		return nil, err
	}
	em, err := errorMetric.CurryWith(map[string]string{DefaultLabels.Service: cfg.ServiceID})
	if err != nil {
		return nil, err
	}
	return &ErrorMetric{em}, nil
}
