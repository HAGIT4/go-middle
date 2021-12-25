package api

import (
	"net/http"

	"github.com/HAGIT4/go-middle/internal/server/service"
)

type updateHandler struct {
	metricService *service.MetricService
}

func newUpdateHandler() *updateHandler {
	serv := service.NewMetricService()
	return &updateHandler{
		metricService: serv,
	}
}

func (h updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("application-type", "text/plain")
	if r.Method != http.MethodPost {
		err := newHandlerForbiddenMethodError(r.Method)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
	}

	if h := r.Header.Get("application-type"); h != "text/plain" {
		err := newHandlerApplicationTypeError(h)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
	}

	metricType, err := parseMetricType(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if metricType == metricTypeGauge {
		metricInfo, err := parseMetricGauge(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		err = h.metricService.UpdateGauge(metricInfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		resp := []byte{}
		w.Write(resp)
	} else if metricType == metricTypeCounter {
		metricInfo, err := parseMetricCounter(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		err = h.metricService.UpdateCounter(metricInfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		resp := []byte{}
		w.Write(resp)
	}
}
