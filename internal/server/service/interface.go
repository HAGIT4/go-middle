package service

import "github.com/HAGIT4/go-middle/internal/server/models"

type IMetricService interface {
	UpdateGauge(metricInfo *models.MetricGaugeInfo) error
	UpdateCounter(metricInfo *models.MetricCounterInfo) error
}
