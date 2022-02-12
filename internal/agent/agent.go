package agent

import (
	"net/http"
	"time"

	"github.com/HAGIT4/go-middle/pkg/agent/config"
	"github.com/rs/zerolog"
)

type agent struct {
	serverAddr     string
	pollInterval   time.Duration
	reportInterval time.Duration
	httpClient     *http.Client
	hashKey        string
	batch          bool
	logger         *zerolog.Logger
}

var _ AgentInterface = (*agent)(nil)

func NewAgent(cfg *config.AgentConfig) (a *agent, err error) {
	httpClient := &http.Client{}
	logger, err := NewAgentLogger()
	if err != nil {
		return nil, err
	}
	a = &agent{
		serverAddr:     cfg.ServerAddr,
		pollInterval:   cfg.PollInterval,
		reportInterval: cfg.ReportInterval,
		httpClient:     httpClient,
		hashKey:        cfg.HashKey,
		batch:          cfg.Batch,
		logger:         logger,
	}
	return a, nil
}

func (a *agent) CollectMetrics() *agentData {
	data := newAgentData()
	return data
}
