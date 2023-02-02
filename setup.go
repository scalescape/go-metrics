package metrics

import (
	"fmt"
	"log"

	"github.com/scalescape/go-metrics/common"
	"github.com/scalescape/go-metrics/prometheus"
)

type Kind int

const (
	Prometheus = iota
	InfluxDB
	Statsd
)

type Config struct {
	Address   string `split_words:"true" required:"true"`
	ServiceID string `split_words:"true" default:"service"`
	// const labels
	labels      map[string]string
	enablePprof bool
	Kind
}

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
	Kind:      Prometheus,
}

func Setup(opts ...Option) (Observer, error) {
	cfg := DefaultConfig
	for _, o := range opts {
		o(cfg)
	}
	var reporter reporter
	var err error

	if cfg.Kind == Prometheus {
		pcfg := prometheus.Config{
			Address:     cfg.Address,
			ServiceID:   cfg.ServiceID,
			ConstLabels: map[string]string{},
			LabelNames:  common.DefaultLabels,
			EnablePprof: cfg.enablePprof,
		}
		reporter, err = prometheus.Setup(pcfg)
		if err != nil {
			return Observer{}, fmt.Errorf("unable to setup prometheus: %w", err)
		}
		log.Printf("I! using prometheus reporter")
	}
	observer := Observer{
		Config:   *cfg,
		reporter: reporter,
	}
	return observer, nil
}
