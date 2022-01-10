package service_test

import (
	"math"
	"testing"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/pkg/models"
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
	restoreConfig := &models.RestoreConfig{}
	ms, _ := service.NewMetricService(restoreConfig)
	metricName := "new metric"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.UpdateGauge(metricName, tt.value)
			actualValue, err := ms.GetGauge(metricName)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, actualValue)
		})
	}
}
