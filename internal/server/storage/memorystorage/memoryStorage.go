package memorystorage

import (
	dbModels "github.com/HAGIT4/go-middle/pkg/server/storage/models"
)

type MemoryStorage struct {
	storageGauge   map[string]float64
	storageCounter map[string]int64
}

// var _ StorageInterface = (*MemoryStorage)(nil)

func NewMemoryStorage() (ms *MemoryStorage, err error) {
	stGauge := make(map[string]float64)
	stCounter := make(map[string]int64)
	ms = &MemoryStorage{
		storageGauge:   stGauge,
		storageCounter: stCounter,
	}
	return ms, nil
}

func (st *MemoryStorage) UpdateGauge(metricName string, metricValue float64) (err error) {
	st.storageGauge[metricName] = metricValue
	return nil
}

func (st *MemoryStorage) GetGauge(metricName string) (metricValue float64, err error) {
	metricValue, found := st.storageGauge[metricName]
	if !found {
		err := newStorageGaugeNotFoundError(metricName)
		return 0, err
	}
	return metricValue, nil
}

func (st *MemoryStorage) GetGaugeAll() (metricNameToValue map[string]float64, err error) {
	return st.storageGauge, nil
}

func (st *MemoryStorage) UpdateCounter(metricName string, metricValue int64) (err error) {
	st.storageCounter[metricName] = metricValue
	return nil
}

func (st *MemoryStorage) GetCounter(metricName string) (metricValue int64, err error) {
	metricValue, found := st.storageCounter[metricName]
	if !found {
		err := newStorageCounterNotFoundError(metricName)
		return 0, err
	}
	return metricValue, nil
}

func (st *MemoryStorage) GetCounterAll() (metricNameToValue map[string]int64, err error) {
	return st.storageCounter, nil
}

func (st *MemoryStorage) Ping() (err error) {
	return nil
}

func (st *MemoryStorage) UpdateBatch(req *dbModels.BatchUpdateRequest) (err error) {
	return nil // TODO
}
