package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/HAGIT4/go-middle/pkg/models"
)

func (sv *MetricService) GetHashKey() string {
	return sv.hashKey
}

func (sv *MetricService) CheckHash(metric *models.Metrics) (err error) {
	reqHash, err := hex.DecodeString(metric.Hash)
	if err != nil {
		return err
	}
	var localHash []byte
	h := hmac.New(sha256.New, []byte(sv.hashKey))
	switch metric.MType {
	case "gauge":
		h.Write([]byte(fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value)))
		localHash = h.Sum(nil)
	case "counter":
		h.Write([]byte(fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta)))
		localHash = h.Sum(nil)
	default:
		return newServiceMetricTypeUnknownError(metric.MType)
	}
	if !hmac.Equal(reqHash, localHash) {
		return newHashNotMatchingError(hex.EncodeToString(reqHash), hex.EncodeToString(localHash))
	}
	return nil
}

func (sv *MetricService) ComputeHash(metric *models.Metrics) (err error) {
	h := hmac.New(sha256.New, []byte(sv.hashKey))
	switch metric.MType {
	case "gauge":
		h.Write([]byte(fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value)))
		metric.Hash = hex.EncodeToString(h.Sum(nil))
	case "counter":
		h.Write([]byte(fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta)))
		metric.Hash = hex.EncodeToString(h.Sum(nil))
	default:
		return newServiceMetricTypeUnknownError(metric.MType)
	}
	return nil
}
