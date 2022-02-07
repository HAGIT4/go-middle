package service

import (
	"github.com/HAGIT4/go-middle/internal/server/storage"
	"github.com/HAGIT4/go-middle/pkg/models"
)

type MetricService struct {
	storage       storage.StorageInterface
	restoreConfig models.RestoreConfig
}

var _ MetricServiceInterface = (*MetricService)(nil)

func NewMetricService(restoreConfig *models.RestoreConfig) (serv *MetricService, err error) {
	st, err := storage.NewMemoryStorage()
	if err != nil {
		return nil, err
	}

	if restoreConfig.StoreInterval == 0 {
		restoreConfig.SyncWrite = true
	} else {
		restoreConfig.SyncWrite = false
	}

	serv = &MetricService{
		storage:       st,
		restoreConfig: *restoreConfig,
	}

	return serv, nil
}
