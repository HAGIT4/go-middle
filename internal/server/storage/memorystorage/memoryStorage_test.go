package memorystorage_test

import (
	"log"
	"math"
	"testing"

	"github.com/HAGIT4/go-middle/internal/server/storage/memorystorage"
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
	ms, err := memorystorage.NewMemoryStorage()
	if err != nil {
		log.Fatal(err)
	}
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
