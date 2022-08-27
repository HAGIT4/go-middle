// Package config Server configuration structure
package config

import "time"

// APIConfig defines server configuration
type APIConfig struct {
	RestoreConfig *APIRestoreConfig
	ServerAddr    string `env:"ADDRESS"`
	HashKey       string `env:"KEY"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
	CryptoKey     string `env:"CRYPTO_KEY"`
}

// APIRestoreConfig defines restoring from backup policy
type APIRestoreConfig struct {
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	SyncWrite     bool
}
