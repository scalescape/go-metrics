package metrics

import (
	"fmt"

	"github.com/scalescape/go-metrics/common"
	"github.com/scalescape/go-metrics/prometheus"
	"github.com/scalescape/go-metrics/statsd"
)

type Option func(c *common.Config)

func WithAddress(addr string) Option {
	return func(c *common.Config) {
		c.Address = addr
	}
}

func WithServiceName(s string) Option {
	return func(c *common.Config) {
		c.ServiceID = s
	}
}

func WithPprof() Option {
	return func(c *common.Config) {
		c.EnablePprof = true
	}
}

func WithKind(kind common.Kind) Option {
	return func(c *common.Config) {
		c.Kind = kind
	}
}

var DefaultConfig = &common.Config{
	Address:   ":9100",
	ServiceID: "go-service",
	Kind:      common.Prometheus,
}

func Setup(opts ...Option) (Observer, error) {
	cfg := DefaultConfig
	for _, o := range opts {
		o(cfg)
	}
	var reporter reporter
	var err error

	switch cfg.Kind {
	case common.Prometheus:
		reporter, err = prometheus.Setup(*cfg)
		if err != nil {
			return Observer{}, fmt.Errorf("unable to setup prometheus: %w", err)
		}
	case common.Statsd:
		reporter, err = statsd.Setup(*cfg)
		if err != nil {
			return Observer{}, fmt.Errorf("unable to setup influx: %w", err)
		}
	}

	observer := Observer{
		Config:   *cfg,
		reporter: reporter,
	}
	return observer, nil
}
