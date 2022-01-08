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
}

var _ MetricServerInterface = (*metricServer)(nil)

func NewMetricServer(addr string, restoreConfig *models.RestoreConfig) *metricServer {
	s := service.NewMetricService(restoreConfig)
	httpMux := newMetricRouter(s)

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
