package models

import "time"

type RestoreConfig struct {
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	SyncWrite     bool
}
