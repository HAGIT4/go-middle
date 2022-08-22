// Package service Service for metric collector
package service

import (
	"github.com/HAGIT4/go-middle/internal/server/storage"
	"github.com/HAGIT4/go-middle/pkg/server/service/config"
	"github.com/rs/zerolog"
)

type MetricService struct {
	storage       storage.StorageInterface
	restoreConfig *config.MetricServiceRestoreConfig
	hashKey       string
	logger        *zerolog.Logger
}

var _ MetricServiceInterface = (*MetricService)(nil)

func NewMetricService(cfg *config.MetricServiceConfig) (sv *MetricService, err error) {
	if cfg.RestoreConfig != nil {
		if cfg.RestoreConfig.StoreInterval == 0 {
			cfg.RestoreConfig.SyncWrite = true
		} else {
			cfg.RestoreConfig.SyncWrite = false
		}
	}

	logger, err := NewServiceLogger()
	if err != nil {
		return nil, err
	}

	sv = &MetricService{
		storage:       cfg.Storage,
		restoreConfig: cfg.RestoreConfig,
		hashKey:       cfg.HashKey,
		logger:        logger,
	}

	if sv.restoreConfig != nil {
		if err := sv.RestoreDataFromFile(); err != nil {
			return nil, err
		}
	}

	return sv, nil
}
