package storage

type StorageInterface interface {
	GetGauge(metricName string) (metricValue float64, err error)
	GetGaugeAll() (metricNameToValue map[string]float64, err error)
	GetCounter(metricName string) (metricValue int64, err error)
	GetCounterAll() (metricNameToValue map[string]int64, err error)

	UpdateGauge(metricName string, metricValue float64) (err error)
	UpdateCounter(metricName string, metricValue int64) (err error)

	Ping() (err error)
}
