package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Option func(c *Config)

func WithAddress(addr string) Option {
	return func(c *Config) {
		c.Address = addr
	}
}

func WithServiceName(s string) Option {
	return func(c *Config) {
		c.ServiceID = s
	}
}

func WithPprof() Option {
	return func(c *Config) {
		c.enablePprof = true
	}
}

var DefaultConfig = &Config{
	Address:   ":9100",
	ServiceID: "go-service",
	labels:    DefaultLabels,
}

func Setup(opts ...Option) (Observer, error) {
	sm := http.NewServeMux()
	cfg := DefaultConfig
	for _, o := range opts {
		o(cfg)
	}
	sm.Handle("/metrics", promhttp.Handler())

	latencyMetric, err := NewLatencyMetric(*cfg)
	if err != nil {
		return Observer{}, err
	}
	errorMetric, err := NewErrorMetric(*cfg)
	if err != nil {
		return Observer{}, err
	}

	if cfg.enablePprof {
		registerPprof(sm)
	}

	go func() {
		log.Printf("I! starting metrics server at: %s\n", cfg.Address)
		if err := http.ListenAndServe(cfg.Address, sm); err != nil {
			log.Printf("E! error starting metrics server: %v\n", err)
			return
		}
	}()
	observer := Observer{
		Config:        *cfg,
		latencyMetric: latencyMetric,
		errorMetric:   errorMetric,
	}
	return observer, nil
}
