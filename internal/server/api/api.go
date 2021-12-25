package api

const (
	metricTypeGauge   = "gauge"
	metricTypeCounter = "counter"
)

type metricServer struct {
	addr    string
	handler *metricRouter
}

func NewMetricServer(addr string) *metricServer {
	httpMux := newMetricRouter()

	metricServer := &metricServer{
		addr:    addr,
		handler: httpMux,
	}
	return metricServer
}

func (s *metricServer) ListenAndServe() {
	s.handler.mux.Run(s.addr)
}
