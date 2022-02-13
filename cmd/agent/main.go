package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/HAGIT4/go-middle/internal/agent"
)

func main() {
	cfg, err := agent.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	a, err := agent.NewAgent(cfg)
	if err != nil {
		log.Fatal(err)
	}

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go a.CollectMetrics(stopCh)
	go a.SendMetricsWithInterval(agent.TypeJSON, cfg.Batch, stopCh)
	<-stopCh
}
