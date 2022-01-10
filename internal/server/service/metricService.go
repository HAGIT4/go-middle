package service

import (
	"os"

	"github.com/HAGIT4/go-middle/internal/server/storage"
	"github.com/HAGIT4/go-middle/pkg/models"
)

type MetricService struct {
	storage storage.StorageInterface

	restoreConfig models.RestoreConfig
	restoreFile   *os.File
}

var _ MetricServiceInterface = (*MetricService)(nil)

func NewMetricService(restoreConfig *models.RestoreConfig) (serv *MetricService, err error) {
	st := storage.NewMemoryStorage()

	if restoreConfig.StoreInterval == 0 {
		restoreConfig.SyncWrite = true
	} else {
		restoreConfig.SyncWrite = false
	}

	serv = &MetricService{
		storage:       st,
		restoreConfig: *restoreConfig,
	}

	if err := serv.RestoreDataFromFile(); err != nil {
		return nil, err
	}

	fileFlags := os.O_WRONLY | os.O_CREATE
	restoreFile, err := os.OpenFile(serv.restoreConfig.StoreFile, fileFlags, 0666)
	if err != nil {
		return nil, err
	}
	serv.restoreFile = restoreFile

	return serv, nil
}

// func (s *MetricService) WriteMetricToFileSync(metricInfo *models.Metrics) (err error) {
// 	metricBz, err := json.Marshal(metricInfo)
// 	if err != nil {
// 		return err
// 	}
// 	metricBz = append(metricBz, '\n')
// 	if _, err = s.restoreFile.Write(metricBz); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (s *MetricService) CloseDataFile() (err error) {
	return s.restoreFile.Close()
}

func (s *MetricService) GetMetricModelsAll() (allMetrics []models.Metrics, err error) {
	gaugeNameToValue, counterNameToValue, err := s.GetMetricAll()
	if err != nil {
		return nil, err
	}
	for name, value := range gaugeNameToValue {
		metricModel := &models.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		}
		allMetrics = append(allMetrics, *metricModel)
	}
	for name, delta := range counterNameToValue {
		metricModel := &models.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &delta,
		}
		allMetrics = append(allMetrics, *metricModel)
	}
	return allMetrics, nil
}
