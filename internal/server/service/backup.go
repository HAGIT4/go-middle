package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/HAGIT4/go-middle/pkg/models"
)

func (s *MetricService) RestoreDataFromFile() (err error) {
	if !s.restoreConfig.Restore {
		fmt.Println("Not restoring from backup..")
		return nil
	}
	openFlags := os.O_RDONLY | os.O_CREATE
	restoreFile, err := os.OpenFile(s.restoreConfig.StoreFile, openFlags, 0600)
	if err != nil {
		return err
	}
	defer func() {
		err = restoreFile.Close()
	}()

	var metrics []models.Metrics
	scan := bufio.NewScanner(restoreFile)
	for scan.Scan() {
		data := scan.Bytes()
		metric := &models.Metrics{}
		if err := json.Unmarshal(data, &metric); err != nil {
			return err
		}
		metrics = append(metrics, *metric)
	}

	for _, metricIt := range metrics {
		metricInfo := metricIt
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
	}
	return nil
}

func (s *MetricService) SaveDataWithInterval() (err error) {
	if s.restoreConfig.SyncWrite {
		fmt.Println("Sync write mode!")
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

func (s *MetricService) WriteAllMetricsToFile() (err error) {
	openFlags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	backupFile, err := os.OpenFile(s.restoreConfig.StoreFile, openFlags, 0600)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(backupFile)
	allMetrics, err := s.GetMetricModelsAll()
	if err != nil {
		return err
	}
	for _, metricIt := range allMetrics {
		metric := metricIt
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

	return nil
}
