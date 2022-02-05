package agent

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HAGIT4/go-middle/pkg/models"
)

type sendType int

const (
	TypePlain sendType = iota
	TypeJSON
)

func prepareHeaders(st sendType, req *http.Request) (err error) {
	switch st {
	case TypePlain:
		req.Header.Set("application-type", "text/plain")
	case TypeJSON:
		req.Header.Set("Content-Type", "application/json")
	default:
		return newUnknownSendTypeError()
	}
	return nil
}

func prepareURL(st sendType, serverAddress string, metricInfo *models.Metrics) (url string, err error) {
	var urlTemplate string
	switch st {
	case TypePlain:
		switch metricInfo.MType {
		case "gauge":
			urlTemplate = "http://%s/update/%s/%s/%f"
			url = fmt.Sprintf(urlTemplate, serverAddress, metricInfo.MType, metricInfo.ID, *metricInfo.Value)
		case "counter":
			urlTemplate = "http://%s/update/%s/%s/%d"
			url = fmt.Sprintf(urlTemplate, serverAddress, metricInfo.MType, metricInfo.ID, *metricInfo.Delta)
		default:
			return "", newUnknownSendTypeError()
		}
	case TypeJSON:
		urlTemplate = "http://%s/update/"
		url = fmt.Sprintf(urlTemplate, serverAddress)
	default:
		return "", newUnknownSendTypeError()
	}
	return url, nil
}

func prepareData(st sendType, metricInfo *models.Metrics) (data io.Reader, err error) {
	switch st {
	case TypePlain:
		return nil, nil
	case TypeJSON:
		metricInfoBytes, err := json.Marshal(metricInfo)
		if err != nil {
			return nil, err
		}
		dataBuffer := bytes.NewBuffer(metricInfoBytes)
		return dataBuffer, nil
	default:
		return nil, newUnknownSendTypeError()
	}
}

func (a *agent) hashData(metric *models.Metrics) (err error) {
	h := hmac.New(sha256.New, []byte(a.hashKey))
	switch metric.MType {
	case "gauge":
		h.Write([]byte(fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value)))
		metric.Hash = hex.EncodeToString(h.Sum(nil))
	case "counter":
		h.Write([]byte(fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta)))
		metric.Hash = hex.EncodeToString(h.Sum(nil))
	default:
		return newUnknownMetricTypeError(metric.MType)
	}
	return nil
}

func (a *agent) SendMetricsBatch(data *agentData, pollCount int64) (err error) {
	dataSlice := []models.Metrics{}
	for metric, valueIt := range *data.agentDataGauge {
		value := valueIt
		reqMetricInfo := &models.Metrics{
			ID:    metric,
			MType: "gauge",
			Value: &value,
		}
		if a.hashKey != "" {
			if err := a.hashData(reqMetricInfo); err != nil {
				return err
			}
		}
		dataSlice = append(dataSlice, *reqMetricInfo)
	}
	reqMetricInfo := &models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollCount,
	}
	if a.hashKey != "" {
		if err := a.hashData(reqMetricInfo); err != nil {
			return err
		}
	}
	dataSlice = append(dataSlice, *reqMetricInfo)

	dataSliceBytes, err := json.Marshal(dataSlice)
	if err != nil {
		return err
	}
	dataBuffer := bytes.NewBuffer(dataSliceBytes)
	url := fmt.Sprintf("http://%s/updates/", a.serverAddr)
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, dataBuffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (a *agent) SendMetrics(st sendType, data *agentData, pollCount int64) (err error) {
	var reqURL string
	var reqData io.Reader
	for metric, value := range *data.agentDataGauge {
		reqMetricInfo := &models.Metrics{
			ID:    metric,
			MType: "gauge",
			Value: &value,
		}
		if a.hashKey != "" {
			a.hashData(reqMetricInfo)
		}

		reqURL, err = prepareURL(st, a.serverAddr, reqMetricInfo)
		if err != nil {
			return err
		}
		reqData, err = prepareData(st, reqMetricInfo)
		if err != nil {
			return err
		}
		ctx := context.TODO()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, reqData)
		if err != nil {
			return err
		}
		prepareHeaders(st, req)

		resp, err := a.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}

	reqMetricInfo := &models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollCount,
	}
	if a.hashKey != "" {
		a.hashData(reqMetricInfo)
	}

	reqURL, err = prepareURL(st, a.serverAddr, reqMetricInfo)
	if err != nil {
		return err
	}
	reqData, err = prepareData(st, reqMetricInfo)
	if err != nil {
		return err
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, reqData)
	if err != nil {
		return err
	}
	prepareHeaders(st, req)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (a *agent) SendMetricsWithInterval(st sendType, batch bool) (err error) {
	var pollCount int64
	var agentData *agentData

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	tickerPoll := time.NewTicker(a.pollInterval)
	tickerSend := time.NewTicker(a.reportInterval)
	cPoll := tickerPoll.C
	cSend := tickerSend.C
	go func() {
		for {
			select {
			case <-cPoll:
				agentData = a.CollectMetrics()
				pollCount += 1
			case <-cSend:
				if a.batch {
					if err := a.SendMetricsBatch(agentData, pollCount); err != nil {
						log.Println(err.Error())
					}
				} else {
					err := a.SendMetrics(st, agentData, pollCount)
					if err != nil {
						log.Println(err.Error())
					}
				}
				pollCount = 0
			}
		}
	}()
	<-quit
	log.Println("Agent shutdown...")
	return nil
}
