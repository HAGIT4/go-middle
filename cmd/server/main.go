package main

import (
	"log"
	"time"

	"github.com/HAGIT4/go-middle/internal/server/api"
	"github.com/HAGIT4/go-middle/pkg/models"
	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddr    string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func main() {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}
	restoreConfig := &models.RestoreConfig{
		StoreInterval: cfg.StoreInterval,
		StoreFile:     cfg.StoreFile,
		Restore:       cfg.Restore,
	}
	s := api.NewMetricServer(cfg.ServerAddr, restoreConfig)
	s.ListenAndServe()
}
