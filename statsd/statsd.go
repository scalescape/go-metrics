package statsd

import (
	"log"

	"github.com/cactus/go-statsd-client/v5/statsd"
	"github.com/scalescape/go-metrics/common"
)

func Setup(cfg common.Config) (*Reporter, error) {
	ccfg := &statsd.ClientConfig{
		Address:   cfg.Address,
		Prefix:    cfg.Prefix,
		TagFormat: statsd.InfixComma,
	}
	client, err := statsd.NewClientWithConfig(ccfg)
	if err != nil {
		return nil, err
	}
	log.Printf("I! configured statsd reporter, pushing to %s", cfg.Address)
	return &Reporter{
		labels:   cfg.ConstLabels,
		client:   client,
		requests: &Metric{Name: "http_requests"},
		latency:  &Metric{Name: "http_latency_ms"},
		rate:     1.0,
	}, nil
}

type Metric struct {
	Name  string
	Value int64
}
