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
	openFlags := os.O_RDONLY
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

	for _, metric := range metrics {
		if err = s.UpdateMetric(&metric); err != nil {
			return err
		}
	}
	return nil
}

func (s *MetricService) SaveDataWithInterval() (err error) {
	saveTicker := time.NewTicker(s.restoreConfig.StoreInterval)
	saveChan := saveTicker.C
	for range saveChan {
		fmt.Println("Saving metrics..")
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
