package service

import "github.com/HAGIT4/go-middle/pkg/models"

type MetricServiceInterface interface {
	GetGauge(metricName string) (metricValue float64, err error)
	GetCounter(metricName string) (metricValue int64, err error)
	GetMetric(metricInfoReq *models.Metrics) (metricInfoResp *models.Metrics, err error)
	GetMetricAll() (gaugeNameToValue map[string]float64, counterNameToValue map[string]int64, err error)

	UpdateGauge(metricName string, metricValue float64) (err error)
	UpdateCounter(metricName string, metricValue int64) (err error)
	UpdateMetric(metricInfo *models.Metrics) (err error)
}
