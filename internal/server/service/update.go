package service

import (
	"github.com/HAGIT4/go-middle/pkg/models"
)

func (s *MetricService) updateGauge(metricName string, metricValue float64) (err error) {
	if err = s.storage.UpdateGauge(metricName, metricValue); err != nil {
		return err
	}
	return nil
}

func (s *MetricService) updateCounter(metricName string, metricValue int64) (err error) {
	knownValue, err := s.getCounter(metricName)
	if err != nil {
		knownValue = 0
	}
	newValue := knownValue + metricValue
	if err = s.storage.UpdateCounter(metricName, newValue); err != nil {
		return err
	}
	return nil
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
		if err := s.updateGauge(metricName, metricValue); err != nil {
			return err
		}
	case "counter":
		if metricInfo.Delta == nil {
			return newServiceNoDeltaUpdateError(metricName)
		}
		metricDelta := *metricInfo.Delta
		if err := s.updateCounter(metricName, metricDelta); err != nil {
			return err
		}
	default:
		return newServiceMetricTypeUnknownError(metricType)
	}
	if s.restoreConfig != nil {
		if s.restoreConfig.SyncWrite {
			if err := s.WriteAllMetricsToFile(); err != nil {
				return err
			}
		}
	}

	return nil
}
