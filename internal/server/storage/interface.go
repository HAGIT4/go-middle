package storage

import (
	"github.com/HAGIT4/go-middle/pkg/server/storage/models"
)

type StorageInterface interface {
	GetGauge(metricName string) (metricValue float64, err error)
	GetGaugeAll() (metricNameToValue map[string]float64, err error)
	GetCounter(metricName string) (metricValue int64, err error)
	GetCounterAll() (metricNameToValue map[string]int64, err error)

	UpdateGauge(metricName string, metricValue float64) (err error)
	UpdateCounter(metricName string, metricValue int64) (err error)

	// GetMetric(req *models.GetRequest) (resp *models.GetResponse, err error)
	// UpdateMetric(req *models.UpdateRequest) (err error)
	UpdateBatch(req *models.BatchUpdateRequest) (err error)

	Ping() (err error)
}
