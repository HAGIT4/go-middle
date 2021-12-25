package service

import "github.com/HAGIT4/go-middle/internal/server/models"

type IMetricService interface {
	GetGauge(metricName string) (metricValue float64, err error)
	GetCounter(metricName string) (metricValue int64, err error)
	GetMetricAll() (gaugeNameToValue map[string]float64, counterNameToValue map[string]int64, err error)
	UpdateGauge(metricInfo *models.MetricGaugeInfo) (err error)
	UpdateCounter(metricInfo *models.MetricCounterInfo) (err error)
}
