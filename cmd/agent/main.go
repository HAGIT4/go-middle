package main

import (
	"log"
	"time"

	"github.com/HAGIT4/go-middle/internal/agent"
	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddr     string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

func main() {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}
	a := agent.NewAgent(cfg.ServerAddr, cfg.PollInterval, cfg.ReportInterval)
	a.SendMetricsWithInterval()
}
