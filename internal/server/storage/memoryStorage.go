package storage

import (
	"github.com/HAGIT4/go-middle/internal/server/models"
)

type MemoryStorage struct {
	storageGauge   map[string]float64
	storageCounter map[string]int64
}

var _ IStorage = (*MemoryStorage)(nil)

func NewMemoryStorage() *MemoryStorage {
	stGauge := make(map[string]float64)
	stCounter := make(map[string]int64)
	return &MemoryStorage{
		storageGauge:   stGauge,
		storageCounter: stCounter,
	}
}

func (st *MemoryStorage) UpdateGauge(metricInfo *models.MetricGaugeInfo) (err error) {
	metric := metricInfo.Name
	value := metricInfo.Value
	st.storageGauge[metric] = value
	return nil
}

func (st *MemoryStorage) GetGauge(metricName string) (metricValue float64, err error) {
	metricValue, found := st.storageGauge[metricName]
	if !found {
		err := newStorageGaugeNotFoundError(metricName)
		return 0, err // how to return nil value?
	}
	return metricValue, nil
}

func (st *MemoryStorage) UpdateCounter(metricInfo *models.MetricCounterInfo) (err error) {
	metric := metricInfo.Name
	value := metricInfo.Value
	if storedValue, found := st.storageCounter[metric]; !found {
		st.storageCounter[metric] = value
	} else {
		st.storageCounter[metric] = storedValue + value
	}
	return nil
}

func (st *MemoryStorage) GetCounter(metricName string) (metricValue int64, err error) {
	metricValue, found := st.storageCounter[metricName]
	if !found {
		err := newStorageCounterNotFoundError(metricName)
		return 0, err // how to return nil value?
	}
	return metricValue, nil
}
