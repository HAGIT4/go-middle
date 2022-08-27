package api

import (
	"flag"
	"time"

	"github.com/HAGIT4/go-middle/pkg/server/api/config"
	"github.com/caarlos0/env"
)

var (
	addressFlag       *string
	restoreFlag       *bool
	storeIntervalFlag *time.Duration
	storeFileFlag     *string
	hashKeyFlag       *string
	databaseDSNflag   *string
	cryptoKeyFlag     *string
)

func InitConfig() (cfg *config.APIConfig, err error) {
	addressFlag = flag.String("a", "localhost:8080", "Server address:port")
	restoreFlag = flag.Bool("r", true, "True to restore data from file")
	storeIntervalFlag = flag.Duration("i", 300*time.Second, "Backup to file interval")
	storeFileFlag = flag.String("f", "/tmp/devops-metrics-db.json", "File to backup")
	hashKeyFlag = flag.String("k", "", "Key for hashing")
	databaseDSNflag = flag.String("d", "", "Database DSN")
	cryptoKeyFlag = flag.String("crypto-key", "", "Path to file with private key")
	flag.Parse()

	cfg = &config.APIConfig{}
	restoreCfg := &config.APIRestoreConfig{}
	if err := env.Parse(restoreCfg); err != nil {
		return nil, err
	}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if len(cfg.DatabaseDSN) == 0 {
		cfg.DatabaseDSN = *databaseDSNflag
	}

	if len(cfg.DatabaseDSN) == 0 {
		if len(restoreCfg.StoreFile) == 0 {
			restoreCfg.StoreFile = *storeFileFlag
		}

		restoreCfg.Restore = restoreCfg.Restore || *restoreFlag

		if restoreCfg.StoreInterval == 0 {
			restoreCfg.StoreInterval = *storeIntervalFlag
		}

		cfg.RestoreConfig = restoreCfg
	} else {
		cfg.RestoreConfig = nil
	}

	if len(cfg.ServerAddr) == 0 {
		cfg.ServerAddr = *addressFlag
	}

	if len(cfg.HashKey) == 0 {
		cfg.HashKey = *hashKeyFlag
	}

	if len(cfg.CryptoKey) == 0 {
		cfg.CryptoKey = *cryptoKeyFlag
	}
	return cfg, nil

}
