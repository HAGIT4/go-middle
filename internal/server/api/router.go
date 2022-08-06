package api

import (
	_ "github.com/HAGIT4/go-middle/docs"
	"github.com/HAGIT4/go-middle/internal/server/api/middleware"
	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/internal/server/storage"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type metricRouter struct {
	mux     *gin.Engine
	service service.MetricServiceInterface
	storage storage.StorageInterface
}

// @title MetricService API
// @description Service for collecting metrics from agents
// @version 1.0
// @BasePath /
// @host localhost:8080
func newMetricRouter(sv service.MetricServiceInterface, st storage.StorageInterface) (r *metricRouter, err error) {
	mux := gin.Default()
	mux.Use(gzip.Gzip(gzip.DefaultCompression))
	mux.RedirectTrailingSlash = false
	mux.LoadHTMLFiles("web/template/allMetrics.html")

	mux.POST("/update/", parseJSONrequest(), middleware.CheckHashSHA256Middleware(sv),
		updateHandler(sv))
	mux.POST("/update/:metricType/:metricName/:metricValue", parsePlainTextRequest(plainTextParseMethodPost),
		middleware.CheckHashSHA256Middleware(sv), updateHandler(sv))
	mux.POST("/updates/", parseBatchJSONrequest(), middleware.CheckBatchHashSHA256Middleware(sv), updateBatchHandler(sv))

	mux.GET("/value/:metricType/:metricName", parsePlainTextRequest(plainTextParseMethodGet),
		getHandler(sv, getResponseFormatPlain))
	mux.POST("/value/", parseJSONrequest(),
		getHandler(sv, getResponseFormatJSON))
	mux.GET("/", getAllDataHTMLhandler(sv))

	// database
	mux.GET("/ping", databasePingHandler(st))

	//swagger
	mux.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r = &metricRouter{
		mux:     mux,
		service: sv,
		storage: st,
	}

	return r, err
}
