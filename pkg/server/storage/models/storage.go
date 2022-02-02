package models

const (
	TypeGauge int = iota
	TypeCounter
)

type UpdateRequest struct {
	MetricType   int
	MetricID     string
	GaugeValue   float64
	CounterDelta int64
}

type BatchUpdateRequest struct {
	Metrics *[]UpdateRequest
}
