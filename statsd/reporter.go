package statsd

import (
	"time"

	"github.com/cactus/go-statsd-client/v5/statsd"
)

type Reporter struct {
	labels   map[string]string
	latency  *Metric
	requests *Metric
	client   statsd.Statter
	rate     float32
}

type ResponseWriter interface {
	Status() int64
}

func (rep *Reporter) CaptureLatency(labels map[string]string, latency time.Duration) error {
	err := rep.timing(latency, labels)
	return err
}

func (rep *Reporter) CaptureRequest(labels map[string]string) error {
	return rep.request(labels)
}

func (rep *Reporter) timing(lat time.Duration, labels map[string]string) error {
	tags := rep.convertLabels(labels)
	return rep.client.Timing(rep.latency.Name, lat.Milliseconds(), rep.rate, tags...)
}

func (rep *Reporter) request(labels map[string]string) error {
	tags := rep.convertLabels(labels)
	return rep.client.Inc(rep.requests.Name, 1, rep.rate, tags...)
}

func (rep *Reporter) convertLabels(labels map[string]string) []statsd.Tag {
	var tags []statsd.Tag
	// add constant labels
	for k, val := range rep.labels {
		tags = append(tags, statsd.Tag{k, val})
	}
	for k, val := range labels {
		tags = append(tags, statsd.Tag{k, val})
	}
	return tags
}
