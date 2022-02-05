package config

import "time"

type ApiConfig struct {
	RestoreConfig *ApiRestoreConfig
	ServerAddr    string `env:"ADDRESS"`
	HashKey       string `env:"KEY"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
}

type ApiRestoreConfig struct {
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	SyncWrite     bool
}
