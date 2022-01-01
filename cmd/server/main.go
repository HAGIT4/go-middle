package main

import "github.com/HAGIT4/go-middle/internal/server/api"

func main() {
	s := api.NewMetricServerV1(":8080")
	s.ListenAndServe()
}
