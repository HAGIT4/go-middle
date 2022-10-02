package agent

import (
	"encoding/json"
	"flag"
	"os"
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
	configFileFlag     *string
	grpcPortFlag       *int
)

func parseJSON(path string) (cfg *config.AgentConfig, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func InitConfig() (cfg *config.AgentConfig, err error) {
	addressFlag = flag.String("a", "localhost:8080", "Server address:port")
	reportIntervalFlag = flag.Duration("r", 10*time.Second, "Metric report interval")
	pollIntervalFlag = flag.Duration("p", 2*time.Second, "Metric poll interval")
	hashKeyFlag = flag.String("k", "", "SHA256 key for hashing")
	batchFlag = flag.Bool("b", false, "True for batch mode")
	cryptoKeyFlag = flag.String("crypto-key", "", "Path to file with public key")
	configFileFlag = flag.String("c", "", "Path to config JSON")
	grpcPortFlag = flag.Int("g", 0, "Grpc port")
	flag.Parse()

	cfg = &config.AgentConfig{}
	envCfg := &config.AgentConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	if len(envCfg.ConfigFile) != 0 {
		jsonCfg, err := parseJSON(envCfg.ConfigFile)
		if err != nil {
			return nil, err
		}
		cfg = jsonCfg
	} else if len(*configFileFlag) != 0 {
		jsonCfg, err := parseJSON(*configFileFlag)
		if err != nil {
			return nil, err
		}
		cfg = jsonCfg
	}

	switch {
	case len(envCfg.ServerAddr) != 0:
		cfg.ServerAddr = envCfg.ServerAddr
	case len(*addressFlag) != 0:
		cfg.ServerAddr = *addressFlag
	}

	switch {
	case envCfg.ReportInterval != 0*time.Second:
		cfg.ReportInterval = envCfg.ReportInterval
	case *reportIntervalFlag != 0*time.Second:
		cfg.ReportInterval = *reportIntervalFlag
	}

	switch {
	case envCfg.PollInterval != 0*time.Second:
		cfg.PollInterval = envCfg.PollInterval
	case *pollIntervalFlag != 0*time.Second:
		cfg.PollInterval = *pollIntervalFlag
	}

	switch {
	case len(envCfg.HashKey) != 0:
		cfg.HashKey = envCfg.HashKey
	case len(*hashKeyFlag) != 0:
		cfg.HashKey = *hashKeyFlag
	}

	switch {
	case len(envCfg.HashKey) != 0:
		cfg.HashKey = envCfg.HashKey
	case len(*hashKeyFlag) != 0:
		cfg.HashKey = *hashKeyFlag
	}

	switch {
	case envCfg.Batch:
		cfg.Batch = envCfg.Batch
	case *batchFlag:
		cfg.Batch = *batchFlag
	}

	switch {
	case len(envCfg.CryptoKey) != 0:
		cfg.CryptoKey = envCfg.CryptoKey
	case len(*cryptoKeyFlag) != 0:
		cfg.CryptoKey = *cryptoKeyFlag
	}

	switch {
	case envCfg.GrpcPort != 0:
		cfg.GrpcPort = envCfg.GrpcPort
	case *grpcPortFlag != 0:
		cfg.GrpcPort = *grpcPortFlag
	}

	return cfg, nil
}
