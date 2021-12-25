package storage

import "github.com/HAGIT4/go-middle/internal/server/models"

type IStorage interface {
	GetGauge(metricName string) (metricValue float64, err error)
	GetGaugeAll() (metricNameToValue map[string]float64, err error)
	GetCounter(metricName string) (metricValue int64, err error)
	GetCounterAll() (metricNameToValue map[string]int64, err error)

	UpdateGauge(metricInfo *models.MetricGaugeInfo) error
	UpdateCounter(metricInfo *models.MetricCounterInfo) error
}
