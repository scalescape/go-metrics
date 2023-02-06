package metrics

import (
	"log"
	"net/http"
	"time"

	"github.com/scalescape/go-metrics/common"
)

type reporter interface {
	CaptureLatency(map[string]string, time.Duration) error
	CaptureRequest(map[string]string) error
}

type Observer struct {
	common.Config
	reporter
}

func (o Observer) Middleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		mw := common.NewResponseWriter(w)
		next.ServeHTTP(mw, r)
		latency := time.Since(start)

		if err := o.CaptureLatency(common.LatencyLabels(r, mw), latency); err != nil {
			log.Printf("E! error collecting latency metric: %v", err)
			return
		}

		if err := o.CaptureRequest(common.ErrorLabels(r, mw)); err != nil {
			log.Printf("E! error collecting count metric: %v", err)
			return
		}
	}
	return http.HandlerFunc(f)
}
