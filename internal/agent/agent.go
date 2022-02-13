package agent

import (
	"net/http"
	"os"
	"sync"
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

	dataBuffer *agentData
	mu         sync.Mutex
	tickerPoll *time.Ticker
	tickerSend *time.Ticker
	pollCount  int64
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
		mu:             sync.Mutex{},
		tickerPoll:     time.NewTicker(cfg.PollInterval),
		tickerSend:     time.NewTicker(cfg.ReportInterval),
	}
	return a, nil
}

func (a *agent) CollectMetrics(stopCh <-chan os.Signal) (err error) {
	cPoll := a.tickerPoll.C
Loop:
	for {
		select {
		case <-cPoll:
			a.mu.Lock()
			a.dataBuffer, err = newAgentData()
			if err != nil {
				return err
			}
			a.pollCount += 1
			a.mu.Unlock()
		case <-stopCh:
			break Loop
		}
	}
	return nil
}
