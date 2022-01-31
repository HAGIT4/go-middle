package service

import (
	"github.com/HAGIT4/go-middle/pkg/models"
)

func (sv *MetricService) updateGauge(metricName string, metricValue float64) (err error) {
	if err = sv.storage.UpdateGauge(metricName, metricValue); err != nil {
		return err
	}
	return nil
}

func (sv *MetricService) updateCounter(metricName string, metricValue int64) (err error) {
	knownValue, err := sv.getCounter(metricName)
	if err != nil {
		knownValue = 0
	}
	newValue := knownValue + metricValue
	if err = sv.storage.UpdateCounter(metricName, newValue); err != nil {
		return err
	}
	return nil
}

func (sv *MetricService) UpdateMetric(metricInfo *models.Metrics) (err error) {
	metricType := metricInfo.MType
	metricName := metricInfo.ID
	switch metricType {
	case "gauge":
		if metricInfo.Value == nil {
			return newServiceNoValueUpdateError(metricName)
		}
		metricValue := *metricInfo.Value
		if err := sv.updateGauge(metricName, metricValue); err != nil {
			return err
		}
	case "counter":
		if metricInfo.Delta == nil {
			return newServiceNoDeltaUpdateError(metricName)
		}
		metricDelta := *metricInfo.Delta
		if err := sv.updateCounter(metricName, metricDelta); err != nil {
			return err
		}
	default:
		return newServiceMetricTypeUnknownError(metricType)
	}
	if sv.restoreConfig != nil {
		if sv.restoreConfig.SyncWrite {
			if err := sv.WriteAllMetricsToFile(); err != nil {
				return err
			}
		}
	}

	return nil
}
