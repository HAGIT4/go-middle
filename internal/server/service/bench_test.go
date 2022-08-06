package service_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/internal/server/storage/memorystorage"
	"github.com/HAGIT4/go-middle/pkg/models"
	"github.com/HAGIT4/go-middle/pkg/server/service/config"
)

var (
	ms     *memorystorage.MemoryStorage
	cfg    *config.MetricServiceConfig
	metric *models.Metrics
	sv     *service.MetricService
	err    error
)

func initService() {
	ms, err = memorystorage.NewMemoryStorage()
	if err != nil {
		panic("")
	}
	cfg = &config.MetricServiceConfig{
		RestoreConfig: nil,
		HashKey:       "aaaa",
		Storage:       ms,
	}
	sv, err = service.NewMetricService(cfg)
	if err != nil {
		panic("")
	}

	metric = &models.Metrics{
		ID:    "test",
		MType: "gauge",
		Delta: nil,
		Value: nil,
	}
	value := float64(100.0)
	metric.Value = &value
}

func BenchmarkCheckHashPositive(b *testing.B) {
	initService()
	h := hmac.New(sha256.New, []byte(cfg.HashKey))
	h.Write([]byte(fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value)))
	metric.Hash = hex.EncodeToString(h.Sum(nil))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sv.CheckHash(metric)
	}
}

func BenchmarkUpdateMetric(b *testing.B) {
	initService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sv.UpdateMetric(metric)
	}
}
