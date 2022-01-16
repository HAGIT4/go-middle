package main

import (
	"log"
	"time"

	"github.com/HAGIT4/go-middle/internal/agent"
)

type Config struct {
	ServerAddr     string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

func main() {
	cfg, err := agent.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	a := agent.NewAgent(cfg.ServerAddr, cfg.PollInterval, cfg.ReportInterval)
	a.SendMetricsWithInterval()
}
