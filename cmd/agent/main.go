package main

import (
	"flag"
	"log"
	"time"

	"github.com/HAGIT4/go-middle/internal/agent"
	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddr     string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

var (
	addressFlag        *string
	reportIntervalFlag *time.Duration
	pollIntervalFlag   *time.Duration

	address        string
	reportInterval time.Duration
	pollInterval   time.Duration
)

func init() {
	addressFlag = flag.String("a", "localhost:8080", "Server address:port")
	reportIntervalFlag = flag.Duration("r", 10*time.Second, "Metric report interval")
	pollIntervalFlag = flag.Duration("p", 2*time.Second, "Metric poll interval")
}

func main() {
	flag.Parse()
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}

	if len(cfg.ServerAddr) == 0 {
		address = *addressFlag
	} else {
		address = cfg.ServerAddr
	}

	if cfg.ReportInterval == 0 {
		reportInterval = *reportIntervalFlag
	} else {
		reportInterval = cfg.ReportInterval
	}

	if cfg.PollInterval == 0 {
		pollInterval = *pollIntervalFlag
	} else {
		pollInterval = cfg.PollInterval
	}

	a := agent.NewAgent(address, pollInterval, reportInterval)
	a.SendMetricsWithInterval()
}
