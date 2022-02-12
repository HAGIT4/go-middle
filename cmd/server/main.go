package main

import (
	"log"

	"github.com/HAGIT4/go-middle/internal/server/api"
)

func main() {
	cfg, err := api.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	s, err := api.NewMetricServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err = s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
