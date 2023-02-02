package common

import (
	"fmt"
	"net/http"
)

type Labels struct {
	Service string
	Method  string
	Path    string
	Status  string
}

var DefaultLabels = Labels{
	Service: "service",
	Method:  "method",
	Path:    "path",
	Status:  "status",
}

func LatencyLabels(r *http.Request, mw *ResponseWriter) map[string]string {
	return map[string]string{
		DefaultLabels.Method: r.Method,
		DefaultLabels.Path:   r.URL.Path,
		//TODO: mask it as 4xx,5xx,...
		DefaultLabels.Status: fmt.Sprintf("%d", mw.StatusCode),
	}
}
func ErrorLabels(r *http.Request, mw *ResponseWriter) map[string]string {
	return map[string]string{
		DefaultLabels.Method: r.Method,
		DefaultLabels.Path:   r.URL.Path,
		DefaultLabels.Status: fmt.Sprintf("%d", mw.StatusCode),
	}
}
