package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type LatencyMetric struct {
	prometheus.ObserverVec
}

type ErrorMetric struct {
	*prometheus.GaugeVec
}

type ResponseWriter interface {
	Status() int64
}

func (l *LatencyMetric) Capture(labels map[string]string, latency time.Duration) error {
	metric, err := l.GetMetricWith(labels)
	if err != nil {
		return err
	}
	metric.Observe(float64(latency.Milliseconds()))
	return nil
}

func (m *ErrorMetric) Capture(labels map[string]string) error {
	metric, err := m.GetMetricWith(labels)
	if err != nil {
		return err
	}
	metric.Inc()
	return nil
}

func NewLatencyMetric(cfg Config) (*LatencyMetric, error) {
	labNames := cfg.LabelNames
	latencyMetric := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "http_latency_ms"},
		[]string{labNames.Service, labNames.Method, labNames.Path, labNames.Status},
	)
	if err := prometheus.Register(latencyMetric); err != nil {
		return nil, err
	}
	lm, err := latencyMetric.CurryWith(map[string]string{labNames.Service: cfg.ServiceID})
	if err != nil {
		return nil, err
	}
	return &LatencyMetric{lm}, nil
}

func NewErrorMetric(cfg Config) (*ErrorMetric, error) {
	labNames := cfg.LabelNames
	errorMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "http_requests"},
		[]string{labNames.Service, labNames.Method, labNames.Path, labNames.Status})
	if err := prometheus.Register(errorMetric); err != nil {
		return nil, err
	}
	em, err := errorMetric.CurryWith(map[string]string{labNames.Service: cfg.ServiceID})
	if err != nil {
		return nil, err
	}
	return &ErrorMetric{em}, nil
}
