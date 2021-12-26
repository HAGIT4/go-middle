package models

type GetAllMetricsData struct {
	GaugeNameToValue   map[string]float64
	CounterNameToValue map[string]int64
}
