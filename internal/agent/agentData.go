package agent

import (
	"math/rand"
	"runtime"
)

type agentDataGaugeV1 map[string]int

type agentDataV1 struct {
	*agentDataGaugeV1
}

func newAgentDataV1() *agentDataV1 {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	dataGauge := newAgentDataGaugeV1(memStats)
	data := &agentDataV1{
		agentDataGaugeV1: dataGauge,
	}
	return data
}

func newAgentDataGaugeV1(memStats *runtime.MemStats) *agentDataGaugeV1 {
	randomValue := rand.Int()
	data := &agentDataGaugeV1{
		"Alloc":         int(memStats.Alloc),
		"BuckHashSys":   int(memStats.BuckHashSys),
		"Frees":         int(memStats.Frees),
		"GCCPUFraction": int(memStats.GCCPUFraction),
		"GCSys":         int(memStats.GCSys),
		"HeapAlloc":     int(memStats.HeapAlloc),
		"HeapIdle":      int(memStats.HeapIdle),
		"HeapInuse":     int(memStats.HeapInuse),
		"HeapObjects":   int(memStats.HeapObjects),
		"HeapReleased":  int(memStats.HeapReleased),
		"HeapSys":       int(memStats.HeapSys),
		"LastGC":        int(memStats.LastGC),
		"Lookups":       int(memStats.Lookups),
		"MCacheInuse":   int(memStats.MCacheInuse),
		"MCacheSys":     int(memStats.MCacheSys),
		"MSpanInuse":    int(memStats.MSpanInuse),
		"MSpanSys":      int(memStats.MSpanInuse),
		"Mallocs":       int(memStats.Mallocs),
		"NextGC":        int(memStats.NextGC),
		"NumForcedGC":   int(memStats.NumForcedGC),
		"NumGC":         int(memStats.NumGC),
		"OtherSys":      int(memStats.OtherSys),
		"PauseTotalNs":  int(memStats.PauseTotalNs),
		"StackInuse":    int(memStats.StackInuse),
		"StackSys":      int(memStats.StackSys),
		"Sys":           int(memStats.Sys),

		"RandomValue": randomValue,
	}
	return data
}
