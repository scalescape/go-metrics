package prometheus

import "github.com/scalescape/go-metrics/common"

type Config struct {
	Address   string `split_words:"true" required:"true"`
	ServiceID string `split_words:"true" default:"service"`
	// const labels
	ConstLabels map[string]string
	EnablePprof bool
	LabelNames  common.Labels
}
