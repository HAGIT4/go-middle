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

	a := agent.NewAgent(cfg)
	a.SendMetricsWithInterval(agent.TypeJSON, true)
}
