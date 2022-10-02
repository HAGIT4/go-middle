// Package config Server configuration structure
package config

import "time"

// APIConfig defines server configuration
type APIConfig struct {
	RestoreConfig *APIRestoreConfig
	ServerAddr    string `env:"ADDRESS" json:"address"`
	HashKey       string `env:"KEY" json:"hash_key"`
	DatabaseDSN   string `env:"DATABASE_DSN" json:"database_dsn"`
	CryptoKey     string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFile    string `env:"CONFIG"`
	TrustedSubnet string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	GrpcPort      int    `env:"GRPC_PORT" json:"grpc_port"`
}

// APIRestoreConfig defines restoring from backup policy
type APIRestoreConfig struct {
	StoreInterval time.Duration `env:"STORE_INTERVAL" json:"store_interval"`
	StoreFile     string        `env:"STORE_FILE" json:"store_file"`
	Restore       bool          `env:"RESTORE" json:"restore"`
	SyncWrite     bool
}
