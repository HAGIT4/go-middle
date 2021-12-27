package agent

import (
	"math/rand"
	"runtime"
)

type agentDataGauge map[string]int

type agentData struct {
	*agentDataGauge
}

func newAgentData() *agentData {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	dataGauge := newAgentDataGauge(memStats)
	data := &agentData{
		agentDataGauge: dataGauge,
	}
	return data
}

func newAgentDataGauge(memStats *runtime.MemStats) *agentDataGauge {
	randomValue := rand.Int()
	data := &agentDataGauge{
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
