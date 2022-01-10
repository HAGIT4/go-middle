package api

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/pkg/models"
)

const (
	metricTypeGauge   = "gauge"
	metricTypeCounter = "counter"
)

type metricServer struct {
	addr    string
	handler *metricRouter
	ser     service.MetricServiceInterface
}

var _ MetricServerInterface = (*metricServer)(nil)

func NewMetricServer(addr string, restoreConfig *models.RestoreConfig) *metricServer {
	s, _ := service.NewMetricService(restoreConfig) // TODO: process err
	httpMux := newMetricRouter(s)

	metricServer := &metricServer{
		addr:    addr,
		handler: httpMux,
		ser:     s,
	}
	return metricServer
}

func (s *metricServer) ListenAndServe() {
	go func() {
		if err := s.handler.mux.Run(s.addr); err != nil {
			log.Fatal(err)
		}
	}()
	go s.ser.SaveDataWithInterval()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("Server shutdown...")
}
