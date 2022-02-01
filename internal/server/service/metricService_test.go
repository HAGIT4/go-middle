package service_test

import (
	"log"
	"math"
	"testing"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/internal/server/storage/memorystorage"
	"github.com/HAGIT4/go-middle/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestServiceUpdateGauge(t *testing.T) {
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
	restoreConfig := &models.RestoreConfig{
		StoreInterval: 300,
		StoreFile:     "/tmp/devops-metrics-db.json",
		Restore:       false,
	}
	st, err := memorystorage.NewMemoryStorage()
	if err != nil {
		log.Fatal(err)
	}
	svCfg := &service.MetricServiceConfig{
		RestoreConfig: restoreConfig,
		Storage:       st,
	}
	ms, err := service.NewMetricService(svCfg)
	if err != nil {
		log.Fatal(err)
	}

	metricName := "new metric"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricInfo := &models.Metrics{
				ID:    metricName,
				MType: "gauge",
				Value: &tt.value,
			}
			if err := ms.UpdateMetric(metricInfo); err != nil {
				t.Fatal(err)
			}
			getResp, err := ms.GetMetric(metricInfo)
			if err != nil {
				t.Fatal(err)
			}
			actualValue := *getResp.Value
			assert.Equal(t, tt.want, actualValue)
		})
	}
}

func TestUpdateCounter(t *testing.T) {
	tests := []struct {
		name  string
		value int64
		want  int64
	}{
		{
			name:  "positive value",
			value: 30,
			want:  30,
		},
	}
	restoreConfig := &models.RestoreConfig{
		StoreInterval: 300,
		StoreFile:     "/tmp/devops-metrics-db.json",
		Restore:       false,
	}
	st, err := memorystorage.NewMemoryStorage()
	if err != nil {
		log.Fatal(err)
	}
	svCfg := &service.MetricServiceConfig{
		Storage:       st,
		RestoreConfig: restoreConfig,
	}
	ms, _ := service.NewMetricService(svCfg)
	metricName := "new counter"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricInfo := &models.Metrics{
				ID:    metricName,
				MType: "counter",
				Delta: &tt.value,
			}
			if err := ms.UpdateMetric(metricInfo); err != nil {
				t.Fatal(err)
			}
			getResp, err := ms.GetMetric(metricInfo)
			if err != nil {
				t.Fatal(err)
			}
			actualValue := *getResp.Delta
			assert.Equal(t, tt.want, actualValue)
		})
	}
}
