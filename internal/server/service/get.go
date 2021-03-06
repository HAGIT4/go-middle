package service

import (
	"github.com/HAGIT4/go-middle/pkg/models"
)

func (sv *MetricService) getGauge(metricName string) (metricValue float64, err error) {
	metricValue, err = sv.storage.GetGauge(metricName)
	if err != nil {
		return 0, err
	}
	return metricValue, nil
}

func (sv *MetricService) getCounter(metricName string) (metricValue int64, err error) {
	metricValue, err = sv.storage.GetCounter(metricName)
	if err != nil {
		return 0, err
	}
	return metricValue, nil
}

func (sv *MetricService) GetMetric(metricInfoReq *models.Metrics) (metricInfoResp *models.Metrics, err error) {
	metricType := metricInfoReq.MType
	metricName := metricInfoReq.ID
	metricInfoResp = &models.Metrics{
		MType: metricType,
		ID:    metricName,
	}
	switch metricType {
	case "gauge":
		metricValue, err := sv.getGauge(metricName)
		if err != nil {
			return nil, err
		}
		metricInfoResp.Value = &metricValue
	case "counter":
		metricDelta, err := sv.getCounter(metricName)
		if err != nil {
			return nil, err
		}
		metricInfoResp.Delta = &metricDelta
	default:
		return nil, newServiceMetricTypeUnknownError(metricType)
	}

	if len(sv.hashKey) > 0 {
		sv.ComputeHash(metricInfoResp)
	}

	return metricInfoResp, nil
}

func (sv *MetricService) GetMetricAll() (gaugeNameToValue map[string]float64, counterNameToValue map[string]int64, err error) {
	gaugeNameToValue, err = sv.storage.GetGaugeAll()
	if err != nil {
		return nil, nil, err
	}
	counterNameToValue, err = sv.storage.GetCounterAll()
	if err != nil {
		return nil, nil, err
	}
	return gaugeNameToValue, counterNameToValue, nil
}

func (sv *MetricService) GetMetricModelsAll() (allMetrics []models.Metrics, err error) {
	gaugeNameToValue, counterNameToValue, err := sv.GetMetricAll()
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
