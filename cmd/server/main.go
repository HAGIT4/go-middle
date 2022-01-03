package main

import (
	"log"

	"github.com/HAGIT4/go-middle/internal/server/api"
	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddr string `env:"ADDRESS" envDefault:"localhost:8080"`
}

func main() {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}

	s := api.NewMetricServer(cfg.ServerAddr)
	s.ListenAndServe()
}
