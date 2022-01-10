package api

type MetricServerInterface interface {
	ListenAndServe() (err error)
}
