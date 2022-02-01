package api

import (
	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/internal/server/storage"
	"github.com/gin-gonic/gin"
)

type metricRouter struct {
	mux     *gin.Engine
	service service.MetricServiceInterface
	storage storage.StorageInterface
}

func newMetricRouter(sv service.MetricServiceInterface, st storage.StorageInterface) (r *metricRouter, err error) {
	mux := gin.Default()
	// mux.Use(middleware.GzipReadMiddleware())
	// mux.Use(middleware.GzipWriteMiddleware())
	mux.RedirectTrailingSlash = false
	mux.LoadHTMLFiles("web/template/allMetrics.html")

	mux.POST("/update/", parseJSONrequest(),
		updateHandler(sv))
	mux.POST("/update/:metricType/:metricName/:metricValue", parsePlainTextRequest(plainTextParseMethodPost),
		updateHandler(sv))

	mux.GET("/value/:metricType/:metricName", parsePlainTextRequest(plainTextParseMethodGet),
		getHandler(sv, getResponseFormatPlain))
	mux.POST("/value/", parseJSONrequest(),
		getHandler(sv, getResponseFormatJSON))
	mux.GET("/", getAllDataHTMLhandler(sv))

	// mux.POST("/update/", parseJSONrequest(), middleware.CheckHashSHA256Middleware(sv),
	// 	updateHandler(sv))
	// mux.POST("/update/:metricType/:metricName/:metricValue", parsePlainTextRequest(plainTextParseMethodPost),
	// 	middleware.CheckHashSHA256Middleware(sv), updateHandler(sv))

	// mux.GET("/value/:metricType/:metricName", parsePlainTextRequest(plainTextParseMethodGet),
	// 	getHandler(sv, getResponseFormatPlain))
	// mux.POST("/value/", parseJSONrequest(),
	// 	getHandler(sv, getResponseFormatJSON))
	// mux.GET("/", getAllDataHTMLhandler(sv))

	// database
	mux.GET("/ping", databasePingHandler(st))

	r = &metricRouter{
		mux:     mux,
		service: sv,
		storage: st,
	}

	return r, err
}
