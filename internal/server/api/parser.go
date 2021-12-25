package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/HAGIT4/go-middle/internal/server/models"
)

var metricToType map[string]string = map[string]string{
	"Alloc":         "gauge",
	"BuckHashSys":   "gauge",
	"Frees":         "gauge",
	"GCCPUFraction": "gauge",
	"GCSys":         "gauge",
	"HeapAlloc":     "gauge",
	"HeapIdle":      "gauge",
	"HeapInuse":     "gauge",
	"HeapObjects":   "gauge",
	"HeapReleased":  "gauge",
	"HeapSys":       "gauge",
	"LastGC":        "gauge",
	"Lookups":       "gauge",
	"MCacheInuse":   "gauge",
	"MCacheSys":     "gauge",
	"MSpanInuse":    "gauge",
	"MSpanSys":      "gauge",
	"Mallocs":       "gauge",
	"NextGC":        "gauge",
	"NumForcedGC":   "gauge",
	"NumGC":         "gauge",
	"OtherSys":      "gauge",
	"PauseTotalNs":  "gauge",
	"StackInuse":    "gauge",
	"StackSys":      "gauge",
	"Sys":           "gauge",
	"PollCount":     "counter",
}

func parseMetricType(r *http.Request) (metricType string, err error) {
	path := r.URL.Path
	pathArgs := strings.Split(path, "/")[1:]
	metricType = pathArgs[1]
	if metricType != metricTypeGauge && metricType != metricTypeCounter {
		err := newUpdateUnknownMetricTypeError(metricType)
		return "", err
	}
	return metricType, nil
}

func parseMetricGauge(r *http.Request) (metricInfo *models.MetricGaugeInfo, err error) {
	path := r.URL.Path
	pathArgs := strings.Split(path, "/")[1:]

	if pathArgs[0] != "update" {
		return &models.MetricGaugeInfo{}, newUpdateInvalidPathFormatError()
	}

	if len(pathArgs) != 4 {
		return &models.MetricGaugeInfo{}, newUpdateInvalidPathFormatError()
	}

	metricName := pathArgs[2]
	if _, found := metricToType[metricName]; !found {
		return &models.MetricGaugeInfo{}, newUpdateUnknownMetricNameError(metricName)
	}
	if metricType := metricToType[metricName]; metricType != metricTypeGauge {
		return &models.MetricGaugeInfo{}, newUpdateUnknownMetricTypeError(metricType)
	}

	metricValue, err := strconv.ParseFloat(pathArgs[3], 64)
	if err != nil {
		return &models.MetricGaugeInfo{}, err
	}
	metricInfo = &models.MetricGaugeInfo{
		Name:  metricName,
		Value: metricValue,
	}
	return metricInfo, nil
}

func parseMetricCounter(r *http.Request) (metricInfo *models.MetricCounterInfo, err error) {
	path := r.URL.Path
	pathArgs := strings.Split(path, "/")[1:]

	if pathArgs[0] != "update" {
		return &models.MetricCounterInfo{}, newUpdateInvalidPathFormatError()
	}

	if len(pathArgs) != 4 {
		return &models.MetricCounterInfo{}, newUpdateInvalidPathFormatError()
	}

	metricName := pathArgs[2]
	if _, found := metricToType[metricName]; !found {
		return &models.MetricCounterInfo{}, newUpdateUnknownMetricNameError(metricName)
	}
	if metricType := metricToType[metricName]; metricType != metricTypeCounter {
		return &models.MetricCounterInfo{}, newUpdateUnknownMetricTypeError(metricType)
	}

	metricValue, err := strconv.ParseInt(pathArgs[3], 10, 64)
	if err != nil {
		return &models.MetricCounterInfo{}, err
	}
	metricInfo = &models.MetricCounterInfo{
		Name:  metricName,
		Value: metricValue,
	}
	return metricInfo, nil
}
