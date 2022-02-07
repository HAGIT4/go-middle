package agent

import (
	"flag"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddr     string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

var (
	addressFlag        *string
	reportIntervalFlag *time.Duration
	pollIntervalFlag   *time.Duration
)

func InitConfig() (cfg *Config, err error) {
	addressFlag = flag.String("a", "localhost:8080", "Server address:port")
	reportIntervalFlag = flag.Duration("r", 10*time.Second, "Metric report interval")
	pollIntervalFlag = flag.Duration("p", 2*time.Second, "Metric poll interval")
	flag.Parse()

	cfg = &Config{}
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

	return cfg, nil
}
