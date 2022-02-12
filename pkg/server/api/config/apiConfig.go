package config

import "time"

type APIConfig struct {
	RestoreConfig *APIRestoreConfig
	ServerAddr    string `env:"ADDRESS"`
	HashKey       string `env:"KEY"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
}

type APIRestoreConfig struct {
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	SyncWrite     bool
}
