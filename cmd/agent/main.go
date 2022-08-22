package main

import (
	"fmt"
	"log"

	"github.com/HAGIT4/go-middle/internal/agent"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func setNA(s *string) {
	if len(*s) == 0 {
		*s = "N/A"
	}
}

func main() {
	setNA(&buildVersion)
	setNA(&buildDate)
	setNA(&buildCommit)
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
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
