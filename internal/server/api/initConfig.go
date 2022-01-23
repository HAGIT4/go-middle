package api

import (
	"flag"
	"time"

	"github.com/HAGIT4/go-middle/pkg/models"
	"github.com/caarlos0/env"
)

type EnvConfig struct {
	ServerAddr    string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

type ServerConfig struct {
	ServerAddr    string
	RestoreConfig *models.RestoreConfig
}

var (
	addressFlag       *string
	restoreFlag       *bool
	storeIntervalFlag *time.Duration
	storeFileFlag     *string
)

func InitConfig() (cfg *ServerConfig, err error) {
	addressFlag = flag.String("a", "localhost:8080", "Server address:port")
	restoreFlag = flag.Bool("r", true, "True to restore data from file")
	storeIntervalFlag = flag.Duration("i", 300*time.Second, "Backup to file interval")
	storeFileFlag = flag.String("f", "/tmp/devops-metrics-db.json", "File to backup")
	flag.Parse()

	envCfg := &EnvConfig{}
	if err := env.Parse(envCfg); err != nil {
		return nil, err
	}

	restoreConfig := &models.RestoreConfig{}
	cfg = &ServerConfig{
		RestoreConfig: restoreConfig,
	}
	if len(envCfg.ServerAddr) == 0 {
		cfg.ServerAddr = *addressFlag
	} else {
		cfg.ServerAddr = envCfg.ServerAddr
	}

	cfg.RestoreConfig.Restore = envCfg.Restore || *restoreFlag

	if envCfg.StoreInterval == 0 {
		cfg.RestoreConfig.StoreInterval = *storeIntervalFlag
	} else {
		cfg.RestoreConfig.StoreInterval = envCfg.StoreInterval
	}

	if len(envCfg.StoreFile) == 0 {
		cfg.RestoreConfig.StoreFile = *storeFileFlag
	} else {
		cfg.RestoreConfig.StoreFile = envCfg.StoreFile
	}

	return cfg, nil

}