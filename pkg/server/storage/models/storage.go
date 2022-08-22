// Package models Models for database
package models

const (
	TypeGauge int = iota
	TypeCounter
)

// UpdateRequest is a database request to update a metric
type UpdateRequest struct {
	MetricType   int
	MetricID     string
	GaugeValue   float64
	CounterDelta int64
}

// BatchUpdateRequest is a database request to update a batch of metcics
type BatchUpdateRequest struct {
	Metrics *[]UpdateRequest
}

// GetRequest is a database request to get a metric
type GetRequest struct {
	MetricType int
	MetricID   string
}

// GetResponse is a response to a get metric request to a database
type GetResponse struct {
	MetricType   int
	MetricID     string
	GaugeValue   float64
	CounterDelta int64
}
