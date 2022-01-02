package service

import (
	"github.com/HAGIT4/go-middle/internal/server/storage"
	"github.com/HAGIT4/go-middle/pkg/models"
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

func (s *MetricService) GetGauge(metricName string) (metricValue float64, err error) {
	metricValue, err = s.storage.GetGauge(metricName)
	if err != nil {
		return 0, err
	}
	return metricValue, nil
}

func (s *MetricService) GetCounter(metricName string) (metricValue int64, err error) {
	metricValue, err = s.storage.GetCounter(metricName)
	if err != nil {
		return 0, err
	}
	return metricValue, nil
}

func (s *MetricService) GetMetric(metricInfoReq *models.Metrics) (metricInfoResp *models.Metrics, err error) {
	metricType := metricInfoReq.MType
	metricName := metricInfoReq.ID
	metricInfoResp = &models.Metrics{
		MType: metricType,
		ID:    metricName,
	}
	switch metricType {
	case "gauge":
		metricValue, err := s.GetGauge(metricName)
		if err != nil {
			return nil, err
		}
		metricInfoResp.Value = &metricValue
	case "counter":
		metricDelta, err := s.GetCounter(metricName)
		if err != nil {
			return nil, err
		}
		metricInfoResp.Delta = &metricDelta
	default:
		return nil, newServiceMetricTypeUnknownError(metricType)
	}
	return metricInfoResp, nil
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

func (s *MetricService) UpdateGauge(metricName string, metricValue float64) (err error) {
	err = s.storage.UpdateGauge(metricName, metricValue)
	return err
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

func (s *MetricService) UpdateMetric(metricInfo *models.Metrics) (err error) {
	metricType := metricInfo.MType
	metricName := metricInfo.ID
	switch metricType {
	case "gauge":
		if metricInfo.Value == nil {
			return newServiceNoValueUpdateError(metricName)
		}
		metricValue := *metricInfo.Value
		if err := s.UpdateGauge(metricName, metricValue); err != nil {
			return err
		}
	case "counter":
		if metricInfo.Delta == nil {
			return newServiceNoDeltaUpdateError(metricName)
		}
		metricDelta := *metricInfo.Delta
		if err := s.UpdateCounter(metricName, metricDelta); err != nil {
			return err
		}
	default:
		return newServiceMetricTypeUnknownError(metricType)
	}
	return nil
}
