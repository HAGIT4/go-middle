package storage

import "github.com/HAGIT4/go-middle/internal/server/models"

type IStorage interface {
	UpdateGauge(metricInfo *models.MetricGaugeInfo) error
	GetGauge(meticName string) (metricValue float64, err error)
	UpdateCounter(metricInfo *models.MetricCounterInfo) error
	GetCounter(meticName string) (metricValue int64, err error)
}
