// Package api Server for metric collector
package api

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/internal/server/storage"
	"github.com/HAGIT4/go-middle/internal/server/storage/memorystorage"
	"github.com/HAGIT4/go-middle/internal/server/storage/postgresstorage"
	"github.com/HAGIT4/go-middle/pb"
	apiConfig "github.com/HAGIT4/go-middle/pkg/server/api/config"
	serviceConfig "github.com/HAGIT4/go-middle/pkg/server/service/config"
	dbConfig "github.com/HAGIT4/go-middle/pkg/server/storage/config"
	"google.golang.org/grpc"
)

const (
	metricTypeGauge   = "gauge"
	metricTypeCounter = "counter"
)

type metricServer struct {
	addr        string
	grpcPort    int
	handler     *metricRouter
	grpcHandler *metricGrpcHandler
	sv          service.MetricServiceInterface
	restore     bool
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

	var prv *rsa.PrivateKey
	if cfg.CryptoKey != "" {
		prv, err = service.GetPrivateKeyFromPem(cfg.CryptoKey)
		if err != nil {
			return nil, err
		}
	}

	var trustedSubnet *net.IPNet
	if cfg.TrustedSubnet != "" {
		_, trustedSubnet, err = net.ParseCIDR(cfg.TrustedSubnet)
		if err != nil {
			return nil, err
		}
	}

	svCfg := &serviceConfig.MetricServiceConfig{
		Storage:          st,
		RestoreConfig:    serviceRestoreCfg,
		HashKey:          cfg.HashKey,
		CryptoPrivateKey: prv,
		TrustedSubnet:    trustedSubnet,
	}

	sv, err := service.NewMetricService(svCfg)
	if err != nil {
		return nil, err
	}

	httpMux, err := newMetricRouter(sv, st)
	if err != nil {
		return nil, err
	}

	grpcMux, err := newGrpcMetricRouter(sv, st)
	if err != nil {
		return nil, err
	}

	ms = &metricServer{
		addr:        cfg.ServerAddr,
		grpcPort:    cfg.GrpcPort,
		handler:     httpMux,
		grpcHandler: grpcMux,
		sv:          sv,
		restore:     restore,
	}
	return ms, nil
}

func (s *metricServer) ListenAndServe() (err error) {
	go func() {
		if err := s.handler.mux.Run(s.addr); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		grpcAddr := fmt.Sprintf(":%d", s.grpcPort)
		listen, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatal(err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterMetricServiceServer(grpcServer, s.grpcHandler)
		if err := grpcServer.Serve(listen); err != nil {
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
