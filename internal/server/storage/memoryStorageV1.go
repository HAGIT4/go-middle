package storage

type MemoryStorageV1 struct {
	storageGauge   map[string]float64
	storageCounter map[string]int64
}

var _ StorageInterfaceV1 = (*MemoryStorageV1)(nil)

func NewMemoryStorageV1() *MemoryStorageV1 {
	stGauge := make(map[string]float64)
	stCounter := make(map[string]int64)
	return &MemoryStorageV1{
		storageGauge:   stGauge,
		storageCounter: stCounter,
	}
}

func (st *MemoryStorageV1) UpdateGauge(metricName string, metricValue float64) (err error) {
	st.storageGauge[metricName] = metricValue
	return nil
}

func (st *MemoryStorageV1) GetGauge(metricName string) (metricValue float64, err error) {
	metricValue, found := st.storageGauge[metricName]
	if !found {
		err := newStorageGaugeNotFoundError(metricName)
		return 0, err
	}
	return metricValue, nil
}

func (st *MemoryStorageV1) GetGaugeAll() (metricNameToValue map[string]float64, err error) {
	return st.storageGauge, nil
}

func (st *MemoryStorageV1) UpdateCounter(metricName string, metricValue int64) (err error) {
	st.storageCounter[metricName] = metricValue
	return nil
}

func (st *MemoryStorageV1) GetCounter(metricName string) (metricValue int64, err error) {
	metricValue, found := st.storageCounter[metricName]
	if !found {
		err := newStorageCounterNotFoundError(metricName)
		return 0, err
	}
	return metricValue, nil
}

func (st *MemoryStorageV1) GetCounterAll() (metricNameToValue map[string]int64, err error) {
	return st.storageCounter, nil
}
