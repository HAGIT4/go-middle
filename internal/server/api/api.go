package api

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	metricTypeGauge   = "gauge"
	metricTypeCounter = "counter"
)

type metricServerV1 struct {
	addr    string
	handler *metricRouterV1
}

func NewMetricServerV1(addr string) *metricServerV1 {
	httpMux := newMetricRouterV1()

	metricServer := &metricServerV1{
		addr:    addr,
		handler: httpMux,
	}
	return metricServer
}

func (s *metricServerV1) ListenAndServe() {
	go func() {
		if err := s.handler.mux.Run(s.addr); err != nil {
			log.Fatal(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("Server shutdown...")
}
