package config

import "time"

type AgentConfig struct {
	ServerAddr     string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"REPORT_INTERVAL"`
	ReportInterval time.Duration `env:"POLL_INTERVAL"`
	HashKey        string        `env:"AGENT_HASH"`
	Batch          bool          `env:"AGENT_BATCH"`
}
