package agent

import (
	"net/http"
	"time"

	"github.com/HAGIT4/go-middle/pkg/agent/config"
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

func NewAgent(cfg *config.AgentConfig) (a *agent) {
	httpClient := &http.Client{}
	a = &agent{
		serverAddr:     cfg.ServerAddr,
		pollInterval:   cfg.PollInterval,
		reportInterval: cfg.ReportInterval,
		httpClient:     httpClient,
		hashKey:        cfg.HashKey,
		batch:          cfg.Batch,
	}
	return a
}

func (a *agent) CollectMetrics() *agentData {
	data := newAgentData()
	return data
}
