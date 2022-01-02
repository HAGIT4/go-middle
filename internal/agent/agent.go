package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HAGIT4/go-middle/pkg/models"
)

type agent struct {
	serverAddr     string
	pollInterval   time.Duration
	reportInterval time.Duration
	httpClient     *http.Client
}

var _ AgentInterface = (*agent)(nil)

func NewAgent(serverAddr string, pollInterval time.Duration, reportInterval time.Duration) (a *agent) {
	httpClient := &http.Client{}
	a = &agent{
		serverAddr:     serverAddr,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		httpClient:     httpClient,
	}
	return a
}

func (a *agent) CollectMetrics() *agentData {
	data := newAgentData()
	return data
}

func (a *agent) SendMetrics(data *agentData, pollCount int64) (err error) {
	urlTemplate := "http://%s/update/%s/%s/%f"
	dataGauge := *data.agentDataGauge

	for metric, value := range dataGauge {
		url := fmt.Sprintf(urlTemplate, a.serverAddr, "gauge", metric, value)
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("application-type", "text/plain")
		resp, err := a.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}

	urlPollCount := fmt.Sprintf("http://%s/update/counter/PollCount/%d", a.serverAddr, pollCount)
	req, err := http.NewRequest(http.MethodPost, urlPollCount, nil)
	if err != nil {
		return err
	}
	req.Header.Set("application-type", "text/plain")
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (a *agent) SendMetricsJSON(data *agentData, pollCount int64) (err error) {
	url := fmt.Sprintf("http://%s/update/", a.serverAddr)
	dataGauge := *data.agentDataGauge

	for metric, value := range dataGauge {
		reqMetricMsg := &models.Metrics{
			ID:    metric,
			MType: "gauge",
			Value: &value,
		}
		reqMetricBytes, err := json.Marshal(reqMetricMsg)
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(reqMetricBytes)
		req, err := http.NewRequest(http.MethodPost, url, buf)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := a.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}
	reqMetricMsg := &models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollCount,
	}
	reqMetricBytes, err := json.Marshal(reqMetricMsg)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(reqMetricBytes)
	req, err := http.NewRequest(http.MethodPost, url, buf)
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

func (a *agent) SendMetricsWithInterval() (err error) {
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
				err := a.SendMetricsJSON(agentData, pollCount)
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
