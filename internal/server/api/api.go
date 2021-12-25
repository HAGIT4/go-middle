package api

import (
	"net/http"
)

const (
	metricTypeGauge   = "gauge"
	metricTypeCounter = "counter"
)

func NewMetricServer(addr string) *http.Server {
	httpMux := newServeMux()

	httpServer := &http.Server{
		Addr:    addr,
		Handler: httpMux,
	}
	return httpServer
}
