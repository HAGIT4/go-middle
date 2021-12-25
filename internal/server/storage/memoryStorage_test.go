package storage_test

import (
	"math"
	"testing"

	"github.com/HAGIT4/go-middle/internal/server/models"
	"github.com/HAGIT4/go-middle/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdateGauge(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  float64
	}{
		{
			name:  "positive value",
			value: 60.0,
			want:  60.0,
		},
		{
			name:  "negative value",
			value: -30.0,
			want:  -30.0,
		},
		{
			name:  "big value",
			value: math.MaxFloat64,
			want:  math.MaxFloat64,
		},
	}
	ms := storage.NewMemoryStorage()
	metricName := "new metric"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricInfo := &models.MetricGaugeInfo{
				Name:  metricName,
				Value: tt.value,
			}
			ms.UpdateGauge(metricInfo)
			actualValue, err := ms.GetGauge(metricName)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, actualValue)
		})
	}
}
