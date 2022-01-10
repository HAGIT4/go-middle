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

func NewMetricServer(addr string, restoreConfig *models.RestoreConfig) (ms *metricServer, err error) {
	s, err := service.NewMetricService(restoreConfig)
	if err != nil {
		return nil, err
	}

	httpMux, err := newMetricRouter(s)
	if err != nil {
		return nil, err
	}

	ms = &metricServer{
		addr:    addr,
		handler: httpMux,
		ser:     s,
	}
	return ms, nil
}

func (s *metricServer) ListenAndServe() (err error) {
	go func() {
		if err := s.handler.mux.Run(s.addr); err != nil {
			log.Fatal(err)
		}
	}()
	go s.ser.SaveDataWithInterval() // err
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("Server shutdown...")
	if err := s.ser.CloseDataFile(); err != nil {
		return err
	}

	return nil
}
