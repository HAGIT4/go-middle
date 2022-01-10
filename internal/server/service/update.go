package service

import "github.com/HAGIT4/go-middle/pkg/models"

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

	// if s.restoreConfig.SyncWrite {
	// 	if err := s.WriteMetricToFileSync(metricInfo); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
