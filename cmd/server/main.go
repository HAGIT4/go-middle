package main

import (
	"fmt"
	"log"

	"net/http"
	_ "net/http/pprof"

	"github.com/HAGIT4/go-middle/internal/server/api"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func setNA(s *string) {
	if len(*s) == 0 {
		*s = "N/A"
	}
}

func main() {
	setNA(&buildVersion)
	setNA(&buildDate)
	setNA(&buildCommit)
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	cfg, err := api.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	s, err := api.NewMetricServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	go http.ListenAndServe(":8081", nil)
	if err = s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
