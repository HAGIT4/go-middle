package agent

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type agentData struct {
	agentDataGauge *map[string]float64
}

func newAgentData() (data *agentData, err error) {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	dataGauge, err := newAgentDataGauge(memStats)
	if err != nil {
		return nil, err
	}
	data = &agentData{
		agentDataGauge: dataGauge,
	}
	return data, nil
}

func newAgentDataGauge(memStats *runtime.MemStats) (data *map[string]float64, err error) {
	randomValue := rand.Float64()
	data = &map[string]float64{
		"Alloc":         float64(memStats.Alloc),
		"BuckHashSys":   float64(memStats.BuckHashSys),
		"Frees":         float64(memStats.Frees),
		"GCCPUFraction": float64(memStats.GCCPUFraction),
		"GCSys":         float64(memStats.GCSys),
		"HeapAlloc":     float64(memStats.HeapAlloc),
		"HeapIdle":      float64(memStats.HeapIdle),
		"HeapInuse":     float64(memStats.HeapInuse),
		"HeapObjects":   float64(memStats.HeapObjects),
		"HeapReleased":  float64(memStats.HeapReleased),
		"HeapSys":       float64(memStats.HeapSys),
		"LastGC":        float64(memStats.LastGC),
		"Lookups":       float64(memStats.Lookups),
		"MCacheInuse":   float64(memStats.MCacheInuse),
		"MCacheSys":     float64(memStats.MCacheSys),
		"MSpanInuse":    float64(memStats.MSpanInuse),
		"MSpanSys":      float64(memStats.MSpanInuse),
		"Mallocs":       float64(memStats.Mallocs),
		"NextGC":        float64(memStats.NextGC),
		"NumForcedGC":   float64(memStats.NumForcedGC),
		"NumGC":         float64(memStats.NumGC),
		"OtherSys":      float64(memStats.OtherSys),
		"PauseTotalNs":  float64(memStats.PauseTotalNs),
		"StackInuse":    float64(memStats.StackInuse),
		"StackSys":      float64(memStats.StackSys),
		"Sys":           float64(memStats.Sys),
		"TotalAlloc":    float64(memStats.TotalAlloc),

		"RandomValue": randomValue,
	}

	ctx := context.Background()
	memStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}
	(*data)["TotalMemory"] = float64(memStat.Total)
	(*data)["FreeMemory"] = float64(memStat.Free)

	cpuStat, err := cpu.PercentWithContext(ctx, 0, true)
	if err != nil {
		return nil, err
	}
	for i, cpu := range cpuStat {
		metricName := fmt.Sprintf("CPUutilization%d", i+1)
		(*data)[metricName] = cpu
	}
	return data, nil
}
