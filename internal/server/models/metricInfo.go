package models

type MetricGaugeInfo struct {
	Name  string
	Value float64
}

type MetricCounterInfo struct {
	Name  string
	Value int64
}
