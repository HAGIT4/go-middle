package agent

import (
	"context"
	"net/http"
	"time"
)

type agent struct {
	serverAddr     string
	pollInterval   time.Duration
	reportInterval time.Duration
	httpClient     *http.Client
	ctx            context.Context
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
