package api

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/HAGIT4/go-middle/pkg/server/api/config"
	"github.com/caarlos0/env"
)

var (
	addressFlag        *string
	restoreFlag        *bool
	storeIntervalFlag  *time.Duration
	storeFileFlag      *string
	hashKeyFlag        *string
	databaseDSNflag    *string
	cryptoKeyFlag      *string
	configFileFlag     *string
	trustedNetworkFlag *string
	grpcPortFlag       *int
)

func parseJSON(path string) (cfg *config.APIConfig, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, cfg); err != nil {
		return nil, err
	}
	var cfgRestore *config.APIRestoreConfig
	if err = json.Unmarshal(b, cfgRestore); err != nil {
		return nil, err
	}
	cfg.RestoreConfig = cfgRestore
	return cfg, nil
}

func InitConfig() (cfg *config.APIConfig, err error) {
	addressFlag = flag.String("a", "localhost:8080", "Server address:port")
	restoreFlag = flag.Bool("r", true, "True to restore data from file")
	storeIntervalFlag = flag.Duration("i", 300*time.Second, "Backup to file interval")
	storeFileFlag = flag.String("f", "/tmp/devops-metrics-db.json", "File to backup")
	hashKeyFlag = flag.String("k", "", "Key for hashing")
	databaseDSNflag = flag.String("d", "", "Database DSN")
	cryptoKeyFlag = flag.String("crypto-key", "", "Path to file with private key")
	configFileFlag = flag.String("c", "", "Path to config JSON")
	trustedNetworkFlag = flag.String("t", "", "Trusted network")
	grpcPortFlag = flag.Int("g", 0, "Grpc port")
	flag.Parse()

	cfg = &config.APIConfig{}
	restoreCfg := &config.APIRestoreConfig{}

	envCfg := &config.APIConfig{}
	envRestoreCfg := &config.APIRestoreConfig{}

	if err := env.Parse(envRestoreCfg); err != nil {
		return nil, err
	}
	if err := env.Parse(envCfg); err != nil {
		return nil, err
	}

	var cfgJSON *config.APIConfig
	if len(envCfg.ConfigFile) != 0 {
		cfgJSON, err = parseJSON(envCfg.ConfigFile)
		if err != nil {
			return nil, err
		}
		cfg = cfgJSON
	} else if len(*configFileFlag) != 0 {
		cfgJSON, err = parseJSON(*configFileFlag)
		if err != nil {
			return nil, err
		}
		cfg = cfgJSON
	}

	switch {
	case len(envCfg.DatabaseDSN) != 0:
		cfg.DatabaseDSN = envCfg.DatabaseDSN
	case len(*databaseDSNflag) != 0:
		cfg.DatabaseDSN = *databaseDSNflag
	}

	if len(cfg.DatabaseDSN) == 0 {
		switch {
		case len(envRestoreCfg.StoreFile) != 0:
			restoreCfg.StoreFile = envRestoreCfg.StoreFile
		case len(*storeFileFlag) != 0:
			restoreCfg.StoreFile = *storeFileFlag
		}

		restoreCfg.Restore = envRestoreCfg.Restore || *restoreFlag

		switch {
		case envRestoreCfg.StoreInterval != 0:
			restoreCfg.StoreInterval = envRestoreCfg.StoreInterval
		case *storeIntervalFlag != 0:
			restoreCfg.StoreInterval = *storeIntervalFlag
		}
		cfg.RestoreConfig = restoreCfg
	} else {
		cfg.RestoreConfig = nil
	}

	switch {
	case len(envCfg.ServerAddr) != 0:
		cfg.ServerAddr = envCfg.ServerAddr
	case len(*addressFlag) != 0:
		cfg.ServerAddr = *addressFlag
	}

	switch {
	case len(envCfg.HashKey) != 0:
		cfg.HashKey = envCfg.HashKey
	case len(*hashKeyFlag) != 0:
		cfg.HashKey = *hashKeyFlag
	}

	switch {
	case len(envCfg.CryptoKey) != 0:
		cfg.CryptoKey = envCfg.CryptoKey
	case len(*cryptoKeyFlag) != 0:
		cfg.CryptoKey = *cryptoKeyFlag
	}

	switch {
	case len(envCfg.TrustedSubnet) != 0:
		cfg.TrustedSubnet = envCfg.TrustedSubnet
	case len(*trustedNetworkFlag) != 0:
		cfg.TrustedSubnet = *trustedNetworkFlag
	}

	switch {
	case envCfg.GrpcPort != 0:
		cfg.GrpcPort = envCfg.GrpcPort
	case *grpcPortFlag != 0:
		cfg.GrpcPort = *grpcPortFlag
	}

	return cfg, nil

}
