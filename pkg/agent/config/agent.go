// Package config Configuration structure for agent
package config

import "time"

// AgentConfig defines agent configuration
type AgentConfig struct {
	ServerAddr     string        `env:"ADDRESS" json:"address"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" json:"poll_interval"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" json:"report_interval"`
	HashKey        string        `env:"AGENT_HASH" json:"hash_key"`
	Batch          bool          `env:"AGENT_BATCH" json:"batch"`
	CryptoKey      string        `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFile     string        `env:"CONFIG"`
	GrpcPort       int           `env:"GRPC_PORT" json:"grcp_port"`
}
