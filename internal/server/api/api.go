package api

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/internal/server/storage"
	"github.com/HAGIT4/go-middle/internal/server/storage/postgresstorage"
	"github.com/HAGIT4/go-middle/pkg/models"
)

const (
	metricTypeGauge   = "gauge"
	metricTypeCounter = "counter"
)

type metricServer struct {
	addr          string
	handler       *metricRouter
	sv            service.MetricServiceInterface
	restoteConfig *models.RestoreConfig
}

var _ MetricServerInterface = (*metricServer)(nil)

func NewMetricServer(addr string, restoreConfig *models.RestoreConfig, hashKey string, databaseDSN string) (ms *metricServer, err error) {
	sv, err := service.NewMetricService(restoreConfig, hashKey)
	if err != nil {
		return nil, err
	}

	var st storage.StorageInterface
	if len(databaseDSN) == 0 {
		st, err = storage.NewMemoryStorage()
		if err != nil {
			return nil, err
		}
	} else {
		postgresCfg := &postgresstorage.PostgresStorageConfig{
			ConnectionString: databaseDSN,
		}
		st, err = postgresstorage.NewPostgresStorage(postgresCfg)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	httpMux, err := newMetricRouter(sv, st)
	if err != nil {
		return nil, err
	}

	ms = &metricServer{
		addr:          addr,
		handler:       httpMux,
		sv:            sv,
		restoteConfig: restoreConfig,
	}
	return ms, nil
}

func (s *metricServer) ListenAndServe() (err error) {
	if s.restoteConfig != nil {
		if err = s.sv.RestoreDataFromFile(); err != nil {
			return err
		}
		go func() {
			if err := s.sv.SaveDataWithInterval(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	go func() {
		if err := s.handler.mux.Run(s.addr); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("Server shutdown...")

	return nil
}
