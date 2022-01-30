package service

import (
	"github.com/HAGIT4/go-middle/pkg/models"
)

func (s *MetricService) getGauge(metricName string) (metricValue float64, err error) {
	metricValue, err = s.storage.GetGauge(metricName)
	if err != nil {
		return 0, err
	}
	return metricValue, nil
}

func (s *MetricService) getCounter(metricName string) (metricValue int64, err error) {
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
		metricValue, err := s.getGauge(metricName)
		if err != nil {
			return nil, err
		}
		metricInfoResp.Value = &metricValue
	case "counter":
		metricDelta, err := s.getCounter(metricName)
		if err != nil {
			return nil, err
		}
		metricInfoResp.Delta = &metricDelta
	default:
		return nil, newServiceMetricTypeUnknownError(metricType)
	}

	if len(s.hashKey) > 0 {
		s.ComputeHash(metricInfoResp)
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

func (s *MetricService) GetMetricModelsAll() (allMetrics []models.Metrics, err error) {
	gaugeNameToValue, counterNameToValue, err := s.GetMetricAll()
	if err != nil {
		return nil, err
	}
	for nameIt, valueIt := range gaugeNameToValue {
		name := nameIt
		value := valueIt
		metricModel := &models.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		}
		allMetrics = append(allMetrics, *metricModel)
	}
	for nameIt, deltaIt := range counterNameToValue {
		name := nameIt
		delta := deltaIt
		metricModel := &models.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &delta,
		}
		allMetrics = append(allMetrics, *metricModel)
	}
	return allMetrics, nil
}
