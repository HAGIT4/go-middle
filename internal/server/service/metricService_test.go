package service_test

import (
	"math"
	"testing"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/internal/server/storage/memorystorage"
	"github.com/HAGIT4/go-middle/pkg/models"
	"github.com/HAGIT4/go-middle/pkg/server/service/config"
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
	restoreConfig := &config.MetricServiceRestoreConfig{
		StoreInterval: 300,
		StoreFile:     "/tmp/devops-metrics-db.json",
		Restore:       false,
	}
	st, err := memorystorage.NewMemoryStorage()
	if err != nil {
		t.Fatal(err)
	}
	svCfg := &config.MetricServiceConfig{
		RestoreConfig: restoreConfig,
		Storage:       st,
	}
	ms, err := service.NewMetricService(svCfg)
	if err != nil {
		t.Fatal(err)
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
	restoreConfig := &config.MetricServiceRestoreConfig{
		StoreInterval: 300,
		StoreFile:     "/tmp/devops-metrics-db.json",
		Restore:       false,
	}
	st, err := memorystorage.NewMemoryStorage()
	if err != nil {
		t.Fatal(err)
	}
	svCfg := &config.MetricServiceConfig{
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

func TestGetMetricAll(t *testing.T) {
	gaugeValue := float64(30.0)
	metricGauge := &models.Metrics{
		ID:    "newGauge",
		MType: "gauge",
		Value: &gaugeValue,
	}
	counterValue := int64(40)
	metricCounter := &models.Metrics{
		ID:    "newCounter",
		MType: "counter",
		Delta: &counterValue,
	}

	restoreConfig := &config.MetricServiceRestoreConfig{
		StoreInterval: 300,
		StoreFile:     "/tmp/devops-metrics-db.json",
		Restore:       false,
	}
	st, err := memorystorage.NewMemoryStorage()
	if err != nil {
		t.Fatal(err)
	}
	svCfg := &config.MetricServiceConfig{
		Storage:       st,
		RestoreConfig: restoreConfig,
	}
	ms, err := service.NewMetricService(svCfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := ms.UpdateMetric(metricGauge); err != nil {
		t.Fatal(err)
	}
	if err := ms.UpdateMetric(metricCounter); err != nil {
		t.Fatal(err)
	}

	var gaugeMap map[string]float64
	var counterMap map[string]int64
	if gaugeMap, counterMap, err = ms.GetMetricAll(); err != nil {
		t.Fatal(err)
	}
	if gaugeMap["newGauge"] != gaugeValue {
		t.Fatal("Gauge value is not as expected")
	}
	if counterMap["newCounter"] != counterValue {
		t.Fatal("Counter value is not as expected")
	}
}

func TestCheckHash(t *testing.T) {
	gaugeValue := float64(40.0)
	hashKey := "xyzi"
	gaugeMetric := &models.Metrics{
		ID:    "newGauge",
		MType: "gauge",
		Value: &gaugeValue,
		Hash:  "5a2fe0c8b3a414875c7ed3a7a119fdda328dc432f5013dc20c5ddf29a3df0531",
	}

	restoreConfig := &config.MetricServiceRestoreConfig{
		StoreInterval: 300,
		StoreFile:     "/tmp/devops-metrics-db.json",
		Restore:       false,
	}
	st, err := memorystorage.NewMemoryStorage()
	if err != nil {
		t.Fatal(err)
	}
	svCfg := &config.MetricServiceConfig{
		Storage:       st,
		RestoreConfig: restoreConfig,
		HashKey:       hashKey,
	}
	ms, err := service.NewMetricService(svCfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := ms.CheckHash(gaugeMetric); err != nil {
		t.Fatal(err)
	}
}
