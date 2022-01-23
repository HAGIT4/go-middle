package agent

import (
	"bytes"
	"context"
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

func prepareUrl(st sendType, serverAddress string, metricInfo *models.Metrics) (url string, err error) {
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

func (a *agent) SendMetrics(st sendType, data *agentData, pollCount int64) (err error) {
	var reqURL string
	var reqData io.Reader
	for metric, value := range *data.agentDataGauge {
		reqMetricInfo := &models.Metrics{
			ID:    metric,
			MType: "gauge",
			Value: &value,
		}
		reqURL, err = prepareUrl(st, a.serverAddr, reqMetricInfo)
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
	reqURL, err = prepareUrl(st, a.serverAddr, reqMetricInfo)
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

	return nil
}

func (a *agent) SendMetricsWithInterval(st sendType) (err error) {
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
				err := a.SendMetrics(st, agentData, pollCount)
				if err != nil {
					log.Println(err.Error())
				}
				pollCount = 0
			}
		}
	}()
	<-quit
	log.Println("Agent shutdown...")
	return nil
}
