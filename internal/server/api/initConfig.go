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
	HashKey       string        `env:"KEY"`
	DatabaseDSN   string        `env:"DATABASE_DSN"`
}

type ServerConfig struct {
	ServerAddr    string
	RestoreConfig *models.RestoreConfig
	HashKey       string
	DatabaseDSN   string
}

var (
	addressFlag       *string
	restoreFlag       *bool
	storeIntervalFlag *time.Duration
	storeFileFlag     *string
	hashKeyFlag       *string
	databaseDSNflag   *string
)

func InitConfig() (cfg *ServerConfig, err error) {
	addressFlag = flag.String("a", "localhost:8080", "Server address:port")
	restoreFlag = flag.Bool("r", true, "True to restore data from file")
	storeIntervalFlag = flag.Duration("i", 300*time.Second, "Backup to file interval")
	storeFileFlag = flag.String("f", "/tmp/devops-metrics-db.json", "File to backup")
	hashKeyFlag = flag.String("k", "", "Key for hashing")
	databaseDSNflag = flag.String("d", "", "Database DSN")
	flag.Parse()

	envCfg := &EnvConfig{}
	if err := env.Parse(envCfg); err != nil {
		return nil, err
	}
	cfg = &ServerConfig{}

	if len(envCfg.DatabaseDSN) == 0 {
		cfg.DatabaseDSN = *databaseDSNflag
	} else {
		cfg.DatabaseDSN = envCfg.DatabaseDSN
	}

	if len(envCfg.DatabaseDSN) == 0 {
		restoreConfig := &models.RestoreConfig{}

		if len(envCfg.StoreFile) == 0 {
			restoreConfig.StoreFile = *storeFileFlag
		} else {
			restoreConfig.StoreFile = envCfg.StoreFile
		}
		restoreConfig.Restore = envCfg.Restore || *restoreFlag

		if envCfg.StoreInterval == 0 {
			restoreConfig.StoreInterval = *storeIntervalFlag
		} else {
			restoreConfig.StoreInterval = envCfg.StoreInterval
		}
		cfg.RestoreConfig = restoreConfig
	} else {
		cfg.RestoreConfig = nil
	}

	if len(envCfg.ServerAddr) == 0 {
		cfg.ServerAddr = *addressFlag
	} else {
		cfg.ServerAddr = envCfg.ServerAddr
	}

	if len(envCfg.HashKey) == 0 {
		cfg.HashKey = *hashKeyFlag
	} else {
		cfg.HashKey = envCfg.HashKey
	}

	return cfg, nil

}
