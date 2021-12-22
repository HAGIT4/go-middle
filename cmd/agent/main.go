package main

import (
	"time"

	"github.com/HAGIT4/go-middle/internal/agent"
)

const (
	pollInterval time.Duration = 2 * time.Second
	serverAddr   string        = "127.0.0.1"
)

func main() {
	a := agent.NewAgent(serverAddr, pollInterval)
	a.SendMetricsWithInterval()
}
