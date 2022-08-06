package config

import (
	"time"

	"github.com/HAGIT4/go-middle/internal/server/storage"
)

// MetricServiceConfig defines a metric service configuration
type MetricServiceConfig struct {
	Storage       storage.StorageInterface
	RestoreConfig *MetricServiceRestoreConfig
	HashKey       string
}

// MetricServiceRestoreConfig defines a metric service restore from backup policy
type MetricServiceRestoreConfig struct {
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	SyncWrite     bool
}
