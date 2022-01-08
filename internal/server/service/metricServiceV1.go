package service

import (
	"github.com/HAGIT4/go-middle/internal/server/storage"
)

type MetricService struct {
	storage storage.StorageInterface
}

var _ MetricServiceInterface = (*MetricService)(nil)

func NewMetricService() *MetricService {
	st := storage.NewMemoryStorage()
	serv := &MetricService{
		storage: st,
	}
	return serv
}

func (s *MetricService) UpdateGauge(metricName string, metricValue float64) (err error) {
	err = s.storage.UpdateGauge(metricName, metricValue)
	return err
}

func (s *MetricService) GetGauge(metricName string) (metricValue float64, err error) {
	metricValue, err = s.storage.GetGauge(metricName)
	if err != nil {
		return 0, err
	}
	return metricValue, nil
}

func (s *MetricService) UpdateCounter(metricName string, metricValue int64) (err error) {
	knownValue, err := s.GetCounter(metricName)
	if err != nil {
		knownValue = 0
	}
	newValue := knownValue + metricValue
	err = s.storage.UpdateCounter(metricName, newValue)
	return err
}

func (s *MetricService) GetCounter(metricName string) (metricValue int64, err error) {
	metricValue, err = s.storage.GetCounter(metricName)
	if err != nil {
		return 0, err
	}
	return metricValue, nil
}

func (s *MetricService) GetMetricAll() (gaugeNameToValue map[string]float64, counterNameToValue map[string]int64, err error) {
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
