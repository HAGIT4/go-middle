package agent

type AgentInterface interface {
	CollectMetrics() (data *agentData)

	SendMetrics(data *agentData, pollCount int64) (err error)
	SendMetricsJSON(data *agentData, pollCount int64) (err error)
	SendMetricsWithInterval() (err error)
}
