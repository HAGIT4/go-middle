package agent

type AgentInterface interface {
	CollectMetrics() (data *agentData)

	SendMetrics(st sendType, data *agentData, pollCount int64) (err error)
}
