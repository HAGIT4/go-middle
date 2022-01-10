package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

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
		fmt.Println(err.Error())
		return nil, err
	}

	fileFlags := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	restoreFile, err := os.OpenFile(serv.restoreConfig.StoreFile, fileFlags, 0222)
	if err != nil {
		return nil, err
	}
	serv.restoreFile = restoreFile

	return serv, nil
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

	if s.restoreConfig.SyncWrite {
		if err := s.WriteMetricToFileSync(metricInfo); err != nil {
			return err
		}
	}
	return nil
}

func (s *MetricService) RestoreDataFromFile() (err error) {
	if s.restoreConfig.Restore {
		openFlags := os.O_RDONLY | os.O_CREATE
		restoreFile, err := os.OpenFile(s.restoreConfig.StoreFile, openFlags, 0444)
		if err != nil {
			return err
		}
		defer func() {
			err = restoreFile.Close()
		}()

		scan := bufio.NewScanner(restoreFile)
		for scan.Scan() {
			data := scan.Bytes()
			metric := &models.Metrics{}
			if err := json.Unmarshal(data, &metric); err != nil {
				return err
			}
			if err := s.UpdateMetric(metric); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *MetricService) WriteMetricToFileSync(metricInfo *models.Metrics) (err error) {
	metricBz, err := json.Marshal(metricInfo)
	if err != nil {
		return err
	}
	metricBz = append(metricBz, '\n')
	if _, err = s.restoreFile.Write(metricBz); err != nil {
		return err
	}
	return nil
}

func (s *MetricService) CloseDataFile() (err error) {
	return s.restoreFile.Close()
}

func (s *MetricService) SaveDataWithInterval() (err error) {
	if !s.restoreConfig.Restore {
		return nil
	}
	saveTicker := time.NewTicker(s.restoreConfig.StoreInterval)
	saveChan := saveTicker.C
	for range saveChan {
		if err = s.WriteAllMetricsToFile(); err != nil {
			return err
		}
	}
	return nil
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

func (s *MetricService) WriteAllMetricsToFile() (err error) {
	tmpFile, err := os.CreateTemp("", "*")
	if err != nil {
		return err
	}
	defer func() {
		err = os.Remove(tmpFile.Name())
	}()

	writer := bufio.NewWriter(tmpFile)
	allMetrics, err := s.GetMetricModelsAll()
	if err != nil {
		return err
	}
	for _, metric := range allMetrics {
		metricBz, err := json.Marshal(metric)
		if err != nil {
			return err
		}
		if _, err := writer.Write(metricBz); err != nil {
			return err
		}
		if err := writer.WriteByte('\n'); err != nil {
			return err
		}
	}
	if err = writer.Flush(); err != nil {
		return err
	}
	if _, err = io.Copy(s.restoreFile, tmpFile); err != nil {
		return err
	}
	return nil
}
