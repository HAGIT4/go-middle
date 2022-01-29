package api

import (
	"github.com/HAGIT4/go-middle/internal/server/api/middleware"
	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/gin-gonic/gin"
)

type metricRouter struct {
	mux     *gin.Engine
	service service.MetricServiceInterface
}

func newMetricRouter(s service.MetricServiceInterface) (r *metricRouter, err error) {
	mux := gin.Default()
	mux.Use(middleware.GzipReadMiddleware())
	mux.Use(middleware.GzipWriteMiddleware())
	mux.RedirectTrailingSlash = false
	mux.LoadHTMLFiles("web/template/allMetrics.html")

	mux.POST("/update/", parseJSONrequest(), updateHandler(s))
	mux.POST("/update/:metricType/:metricName/:metricValue", parsePlainTextRequest(plainTextParseMethodPost), updateHandler(s))
	mux.GET("/value/:metricType/:metricName", parsePlainTextRequest(plainTextParseMethodGet), getHandler(s))
	mux.POST("/value/", parseJSONrequest(), getHandler(s))
	mux.GET("/", getAllDataHTMLhandler(s))

	r = &metricRouter{
		mux:     mux,
		service: s,
	}

	return r, err
}
