package prometheus

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Reporter struct {
	latencyMetric *LatencyMetric
	requests      *ErrorMetric
}

func Setup(cfg Config) (*Reporter, error) {
	latencyMetric, err := NewLatencyMetric(cfg)
	if err != nil {
		return nil, err
	}
	errorMetric, err := NewErrorMetric(cfg)
	if err != nil {
		return nil, err
	}
	reporter := &Reporter{
		latencyMetric: latencyMetric,
		requests:      errorMetric,
	}
	sm := http.NewServeMux()
	if cfg.EnablePprof {
		registerPprof(sm)
	}
	sm.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Printf("I! starting metrics server at: %s\n", cfg.Address)
		if err := http.ListenAndServe(cfg.Address, sm); err != nil {
			log.Printf("E! error starting metrics server: %v\n", err)
			return
		}
	}()
	return reporter, nil
}

func (rep *Reporter) CaptureLatency(labels map[string]string, latency time.Duration) error {
	return rep.latencyMetric.Capture(labels, latency)
}

func (rep *Reporter) CaptureRequest(labels map[string]string) error {
	return rep.requests.Capture(labels)
}

func (rep *Reporter) Close() error {
	//TODO: based on channel signal for stopping
	return nil
}
