package agent

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type agentV1 struct {
	serverAddr     string
	pollInterval   time.Duration
	reportInterval time.Duration
	httpClient     *http.Client
}

func NewAgentV1(serverAddr string, pollInterval time.Duration, reportInterval time.Duration) *agentV1 {
	httpClient := &http.Client{}
	a := &agentV1{
		serverAddr:     serverAddr,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		httpClient:     httpClient,
	}
	return a
}

func (a *agentV1) CollectMetrics() *agentDataV1 {
	data := newAgentDataV1()
	return data
}

func (a *agentV1) SendMetrics(data *agentDataV1, pollCount int64) error {
	urlTemplate := "http://%s/update/%s/%s/%d"
	dataGauge := *data.agentDataGaugeV1

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

func (a *agentV1) SendMetricsWithInterval() error {
	var pollCount int64
	var agentData *agentDataV1

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
				a.SendMetrics(agentData, pollCount)
				pollCount = 0
			}
		}
	}()
	<-quit
	log.Println("Agent shutdown...")
	return nil
}
