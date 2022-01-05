package api

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	metricTypeGauge   = "gauge"
	metricTypeCounter = "counter"
)

type metricServer struct {
	addr          string
	storeInterval time.Duration
	storeFile     string
	restore       bool

	handler *metricRouter
}

var _ MetricServerInterface = (*metricServer)(nil)

func NewMetricServer(addr string, storeInterval time.Duration, storeFile string, restore bool) *metricServer {
	httpMux := newMetricRouter()

	metricServer := &metricServer{
		addr:          addr,
		storeInterval: storeInterval,
		storeFile:     storeFile,
		restore:       restore,

		handler: httpMux,
	}
	return metricServer
}

func (s *metricServer) RestoreDataFromFile() (err error) {
	if !s.restore {
		fmt.Println("Not reading from file..")
		return nil
	}

	storeFile, err := os.Open(s.storeFile)
	if err != nil {
		return err
	}
	defer storeFile.Close()
	// TODO: todo
	return nil
}

func (s *metricServer) SaveWithInterval() {

}

func (s *metricServer) ListenAndServe() {
	s.RestoreDataFromFile()
	go func() {
		if err := s.handler.mux.Run(s.addr); err != nil {
			log.Fatal(err)
		}
	}()
	go s.SaveWithInterval()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("Server shutdown...")
}
