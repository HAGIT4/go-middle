package service

import (
	"fmt"

	"github.com/HAGIT4/go-middle/internal/server/storage"
	"github.com/HAGIT4/go-middle/pkg/models"
)

type MetricServiceConfig struct {
	St            storage.StorageInterface
	RestoreConfig *models.RestoreConfig
	HashKey       string
}

type MetricService struct {
	storage       storage.StorageInterface
	restoreConfig *models.RestoreConfig
	hashKey       string
}

var _ MetricServiceInterface = (*MetricService)(nil)

func NewMetricService(cfg *MetricServiceConfig) (sv *MetricService, err error) {
	if cfg.RestoreConfig != nil {
		if cfg.RestoreConfig.StoreInterval == 0 {
			cfg.RestoreConfig.SyncWrite = true
		} else {
			cfg.RestoreConfig.SyncWrite = false
		}
	}

	sv = &MetricService{
		storage:       cfg.St,
		restoreConfig: cfg.RestoreConfig,
		hashKey:       cfg.HashKey,
	}

	if sv.restoreConfig != nil {
		fmt.Println(sv.restoreConfig)
		if err := sv.RestoreDataFromFile(); err != nil {
			return nil, err
		}
	}

	return sv, nil
}
