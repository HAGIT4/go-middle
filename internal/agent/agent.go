package agent

import (
	"net/http"
	"time"
)

type agent struct {
	serverAddr     string
	pollInterval   time.Duration
	reportInterval time.Duration
	httpClient     *http.Client
	hashKey        string
	batch          bool
}

var _ AgentInterface = (*agent)(nil)

func NewAgent(serverAddr string, pollInterval time.Duration, reportInterval time.Duration, hashKey string, batch bool) (a *agent) {
	httpClient := &http.Client{}
	a = &agent{
		serverAddr:     serverAddr,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		httpClient:     httpClient,
		hashKey:        hashKey,
		batch:          batch,
	}
	return a
}

func (a *agent) CollectMetrics() *agentData {
	data := newAgentData()
	return data
}
