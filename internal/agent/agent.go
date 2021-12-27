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

type agent struct {
	serverAddr     string
	pollInterval   time.Duration
	reportInterval time.Duration
	httpClient     *http.Client
}

func NewAgent(serverAddr string, pollInterval time.Duration, reportInterval time.Duration) *agent {
	httpClient := &http.Client{}
	a := &agent{
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

func (a *agent) SendMetrics(data *agentData, pollCount int64) error {
	urlTemplate := "http://%s/update/%s/%s/%d"
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

func (a *agent) SendMetricsWithInterval() error {
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
				a.SendMetrics(agentData, pollCount)
				pollCount = 0
			}
		}
	}()
	<-quit
	log.Println("Agent shutdown...")
	return nil
}
