package main

import (
	"flag"
	"log"
	"time"

	"github.com/HAGIT4/go-middle/internal/server/api"
	"github.com/HAGIT4/go-middle/pkg/models"
	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddr    string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

var (
	addressFlag       *string
	restoreFlag       *bool
	storeIntervalFlag *time.Duration
	storeFileFlag     *string

	address       string
	restore       bool
	storeInterval time.Duration
	storeFile     string
)

func init() {
	addressFlag = flag.String("a", "localhost:8080", "Server address:port")
	restoreFlag = flag.Bool("r", true, "True to restore data from file")
	storeIntervalFlag = flag.Duration("i", 300*time.Second, "Backup to file interval")
	storeFileFlag = flag.String("f", "/tmp/devops-metrics-db.json", "File to backup")
}

func main() {
	flag.Parse()
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}

	if len(cfg.ServerAddr) == 0 {
		address = *addressFlag
	} else {
		address = cfg.ServerAddr
	}

	restore = cfg.Restore || *restoreFlag

	if cfg.StoreInterval == 0 {
		storeInterval = *storeIntervalFlag
	} else {
		storeInterval = cfg.StoreInterval
	}

	if len(cfg.StoreFile) == 0 {
		storeFile = *storeFileFlag
	} else {
		storeFile = cfg.StoreFile
	}

	restoreConfig := &models.RestoreConfig{
		StoreInterval: storeInterval,
		StoreFile:     storeFile,
		Restore:       restore,
	}

	s, err := api.NewMetricServer(address, restoreConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err = s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
