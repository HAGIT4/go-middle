package service

import "github.com/HAGIT4/go-middle/pkg/models"

type MetricServiceInterface interface {
	// get.go
	getGauge(metricName string) (metricValue float64, err error)
	getCounter(metricName string) (metricValue int64, err error)
	GetMetric(metricInfoReq *models.Metrics) (metricInfoResp *models.Metrics, err error)
	GetMetricAll() (gaugeNameToValue map[string]float64, counterNameToValue map[string]int64, err error)
	GetMetricModelsAll() (allMetrics []models.Metrics, err error)
	// update.go
	updateGauge(metricName string, metricValue float64) (err error)
	updateCounter(metricName string, metricValue int64) (err error)
	UpdateMetric(metricInfo *models.Metrics) (err error)
	UpdateBatch(metricsSlice *[]models.Metrics) (err error)
	// backup.go
	RestoreDataFromFile() (err error)
	SaveDataWithInterval() (err error)
	WriteAllMetricsToFile() (err error)
	// hash.go
	CheckHash(metric *models.Metrics) (err error)
	ComputeHash(metric *models.Metrics) (err error)
	GetHashKey() string
}
