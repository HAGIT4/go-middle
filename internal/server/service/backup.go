package service

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/HAGIT4/go-middle/pkg/models"
)

func (s *MetricService) RestoreDataFromFile() (err error) {
	if s.restoreConfig.Restore {
		openFlags := os.O_RDONLY | os.O_CREATE
		restoreFile, err := os.OpenFile(s.restoreConfig.StoreFile, openFlags, 0666)
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

func (s *MetricService) WriteAllMetricsToFile() (err error) {
	writer := bufio.NewWriter(s.restoreFile)
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

	return nil
}
