package service

type MetricServiceInterface interface {
	GetGauge(metricName string) (metricValue float64, err error)
	GetCounter(metricName string) (metricValue int64, err error)
	GetMetricAll() (gaugeNameToValue map[string]float64, counterNameToValue map[string]int64, err error)
	UpdateGauge(metricName string, metricValue float64) (err error)
	UpdateCounter(metricName string, metricValue int64) (err error)
}
