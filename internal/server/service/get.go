package service

import "github.com/HAGIT4/go-middle/pkg/models"

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
