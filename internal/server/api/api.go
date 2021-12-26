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

type metricServer struct {
	addr    string
	handler *metricRouter
}

func NewMetricServer(addr string) *metricServer {
	httpMux := newMetricRouter()

	metricServer := &metricServer{
		addr:    addr,
		handler: httpMux,
	}
	return metricServer
}

func (s *metricServer) ListenAndServe() {
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
