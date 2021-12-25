package agent

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type agent struct {
	serverAddr   string
	pollInterval time.Duration
	pollCount    int
	httpClient   *http.Client
}

func NewAgent(serverAddr string, pollDuration time.Duration) *agent {
	httpClient := &http.Client{}
	a := &agent{
		serverAddr:   serverAddr,
		pollInterval: pollDuration,
		pollCount:    0,
		httpClient:   httpClient,
	}
	return a
}

func (a *agent) CollectMetrics() *agentData {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	data := newAgentData(memStats, a.pollCount)
	a.pollCount = a.pollCount + 1
	return data
}

func (a *agent) SendMetrics(data *agentData) error {
	urlTemplate := "http://%s/update/%s/%s/%d"
	dataGauge := *data.agentDataGauge
	dataCounter := *data.agentDataCounter

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

	for metric, value := range dataCounter {
		url := fmt.Sprintf(urlTemplate, a.serverAddr, "counter", metric, value)
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

	return nil
}

func (a *agent) CollectAndSendMetrics() error {
	memStats := a.CollectMetrics()
	err := a.SendMetrics(memStats)
	if err != nil {
		return err
	}
	return nil
}

func (a *agent) SendMetricsWithInterval() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	done := make(chan bool, 1)
	go func() {
		<-sigs
		done <- true
	}()

	ticker := time.NewTicker(a.pollInterval)
	c := ticker.C
Loop:
	for {
		select {
		case <-c:
			a.CollectAndSendMetrics()
		case <-done:
			break Loop
		}
	}
	return nil
}
