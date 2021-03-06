package api

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/internal/server/storage"
	"github.com/HAGIT4/go-middle/internal/server/storage/memorystorage"
	"github.com/HAGIT4/go-middle/internal/server/storage/postgresstorage"
	apiConfig "github.com/HAGIT4/go-middle/pkg/server/api/config"
	serviceConfig "github.com/HAGIT4/go-middle/pkg/server/service/config"
	dbConfig "github.com/HAGIT4/go-middle/pkg/server/storage/config"
)

const (
	metricTypeGauge   = "gauge"
	metricTypeCounter = "counter"
)

type metricServer struct {
	addr    string
	handler *metricRouter
	sv      service.MetricServiceInterface
	restore bool
}

var _ MetricServerInterface = (*metricServer)(nil)

func NewMetricServer(cfg *apiConfig.APIConfig) (ms *metricServer, err error) {
	var st storage.StorageInterface
	var postgresCfg = &dbConfig.PostgresStorageConfig{}
	if len(cfg.DatabaseDSN) == 0 {
		st, err = memorystorage.NewMemoryStorage()
		if err != nil {
			return nil, err
		}
	} else {
		postgresCfg.ConnectionString = cfg.DatabaseDSN
		st, err = postgresstorage.NewPostgresStorage(postgresCfg)
		if err != nil {
			return nil, err
		}
	}

	var restore bool
	var serviceRestoreCfg = &serviceConfig.MetricServiceRestoreConfig{}
	if cfg.RestoreConfig != nil {
		serviceRestoreCfg.StoreInterval = cfg.RestoreConfig.StoreInterval
		serviceRestoreCfg.StoreFile = cfg.RestoreConfig.StoreFile
		serviceRestoreCfg.Restore = cfg.RestoreConfig.Restore
		restore = true
	} else {
		serviceRestoreCfg = nil
	}

	svCfg := &serviceConfig.MetricServiceConfig{
		Storage:       st,
		RestoreConfig: serviceRestoreCfg,
		HashKey:       cfg.HashKey,
	}

	sv, err := service.NewMetricService(svCfg)
	if err != nil {
		return nil, err
	}

	httpMux, err := newMetricRouter(sv, st)
	if err != nil {
		return nil, err
	}

	ms = &metricServer{
		addr:    cfg.ServerAddr,
		handler: httpMux,
		sv:      sv,
		restore: restore,
	}
	return ms, nil
}

func (s *metricServer) ListenAndServe() (err error) {
	go func() {
		if err := s.handler.mux.Run(s.addr); err != nil {
			log.Fatal(err)
		}
	}()

	if s.restore {
		go func() {
			if err := s.sv.SaveDataWithInterval(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("Server shutdown...")

	return nil
}
