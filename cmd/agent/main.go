package main

import (
	"log"

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
	a.SendMetricsWithInterval(agent.TypeJSON, true)
}
