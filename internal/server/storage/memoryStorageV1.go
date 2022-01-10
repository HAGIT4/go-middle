package storage

type MemoryStorage struct {
	storageGauge   map[string]float64
	storageCounter map[string]int64
}

var _ StorageInterface = (*MemoryStorage)(nil)

func NewMemoryStorage() *MemoryStorage {
	stGauge := make(map[string]float64)
	stCounter := make(map[string]int64)
	return &MemoryStorage{
		storageGauge:   stGauge,
		storageCounter: stCounter,
	}
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
