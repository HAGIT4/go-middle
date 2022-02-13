package agent

import "os"

type AgentInterface interface {
	CollectMetrics(stopCh <-chan os.Signal) (err error)
	SendMetrics(st sendType, data *agentData, pollCount int64) (err error)
}
