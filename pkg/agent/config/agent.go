// Package config Configuration structure for agent
package config

import "time"

// AgentConfig defines agent configuration
type AgentConfig struct {
	ServerAddr     string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	HashKey        string        `env:"AGENT_HASH"`
	Batch          bool          `env:"AGENT_BATCH"`
}
