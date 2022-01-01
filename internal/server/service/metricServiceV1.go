package service

import (
	"github.com/HAGIT4/go-middle/internal/server/storage"
)

type MetricServiceV1 struct {
	storage storage.StorageInterfaceV1
}

var _ MetricServiceInterfaceV1 = (*MetricServiceV1)(nil)

func NewMetricServiceV1() *MetricServiceV1 {
	st := storage.NewMemoryStorageV1()
	serv := &MetricServiceV1{
		storage: st,
	}
	return serv
}

func (s *MetricServiceV1) UpdateGauge(metricName string, metricValue float64) (err error) {
	err = s.storage.UpdateGauge(metricName, metricValue)
	return err
}

func (s *MetricServiceV1) GetGauge(metricName string) (metricValue float64, err error) {
	metricValue, err = s.storage.GetGauge(metricName)
	if err != nil {
		return 0, err
	}
	return metricValue, nil
}

func (s *MetricServiceV1) UpdateCounter(metricName string, metricValue int64) (err error) {
	knownValue, err := s.GetCounter(metricName)
	if err != nil {
		knownValue = 0
	}
	newValue := knownValue + metricValue
	err = s.storage.UpdateCounter(metricName, newValue)
	return err
}

func (s *MetricServiceV1) GetCounter(metricName string) (metricValue int64, err error) {
	metricValue, err = s.storage.GetCounter(metricName)
	if err != nil {
		return 0, err
	}
	return metricValue, nil
}

func (s *MetricServiceV1) GetMetricAll() (gaugeNameToValue map[string]float64, counterNameToValue map[string]int64, err error) {
	gaugeNameToValue, err = s.storage.GetGaugeAll()
	if err != nil {
		return nil, nil, err
	}
	counterNameToValue, err = s.storage.GetCounterAll()
	if err != nil {
		return nil, nil, err
	}
	return gaugeNameToValue, counterNameToValue, nil
}
