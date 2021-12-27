package service

import (
	"github.com/HAGIT4/go-middle/internal/server/models"
	"github.com/HAGIT4/go-middle/internal/server/storage"
)

type MetricService struct {
	storage storage.IStorage
}

var _ IMetricService = (*MetricService)(nil)

func NewMetricService() *MetricService {
	st := storage.NewMemoryStorage()
	serv := &MetricService{
		storage: st,
	}
	return serv
}

func (s *MetricService) UpdateGauge(metricInfo *models.MetricGaugeInfo) (err error) {
	err = s.storage.UpdateGauge(metricInfo)
	return err
}

func (s *MetricService) GetGauge(metricName string) (metricValue float64, err error) {
	metricValue, err = s.storage.GetGauge(metricName)
	if err != nil {
		return 0, err
	}
	return metricValue, nil
}

func (s *MetricService) UpdateCounter(metricInfo *models.MetricCounterInfo) (err error) {
	knownValue, err := s.GetCounter(metricInfo.Name)
	if err != nil {
		knownValue = 0
	}
	metricInfo = &models.MetricCounterInfo{
		Name:  metricInfo.Name,
		Value: knownValue + metricInfo.Value,
	}
	err = s.storage.UpdateCounter(metricInfo)
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
