package common

type Kind int

const (
	Prometheus = iota
	InfluxDB
	Statsd
)

type Config struct {
	Address   string `split_words:"true" required:"true"`
	Prefix    string
	ServiceID string `split_words:"true" default:"service"`
	// const ConstLabels
	ConstLabels map[string]string
	EnablePprof bool
	Kind
}
