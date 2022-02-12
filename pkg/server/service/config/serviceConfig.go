package config

import (
	"time"

	"github.com/HAGIT4/go-middle/internal/server/storage"
)

type MetricServiceConfig struct {
	Storage       storage.StorageInterface
	RestoreConfig *MetricServiceRestoreConfig
	HashKey       string
}

type MetricServiceRestoreConfig struct {
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	SyncWrite     bool
}
