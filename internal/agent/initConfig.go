package agent

import (
	"flag"
	"time"

	"github.com/HAGIT4/go-middle/pkg/agent/config"
	"github.com/caarlos0/env"
)

var (
	addressFlag        *string
	reportIntervalFlag *time.Duration
	pollIntervalFlag   *time.Duration
	hashKeyFlag        *string
	batchFlag          *bool
	cryptoKeyFlag      *string
)

func InitConfig() (cfg *config.AgentConfig, err error) {
	addressFlag = flag.String("a", "localhost:8080", "Server address:port")
	reportIntervalFlag = flag.Duration("r", 10*time.Second, "Metric report interval")
	pollIntervalFlag = flag.Duration("p", 2*time.Second, "Metric poll interval")
	hashKeyFlag = flag.String("k", "", "SHA256 key for hashing")
	batchFlag = flag.Bool("b", false, "True for batch mode")
	cryptoKeyFlag = flag.String("crypto-key", "", "Path to file with public key")
	flag.Parse()

	cfg = &config.AgentConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if len(cfg.ServerAddr) == 0 {
		cfg.ServerAddr = *addressFlag
	}

	if cfg.ReportInterval == 0*time.Second {
		cfg.ReportInterval = *reportIntervalFlag
	}

	if cfg.PollInterval == 0*time.Second {
		cfg.PollInterval = *pollIntervalFlag
	}

	if cfg.HashKey == "" {
		cfg.HashKey = *hashKeyFlag
	}

	if !cfg.Batch {
		cfg.Batch = *batchFlag
	}

	if len(cfg.CryptoKey) == 0 {
		cfg.CryptoKey = *cryptoKeyFlag
	}

	return cfg, nil
}
