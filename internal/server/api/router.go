package api

import (
	"net/http"
	"strconv"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/gin-gonic/gin"
)

type metricRouterV1 struct {
	mux     *gin.Engine
	service service.MetricServiceInterfaceV1
}

func newMetricRouterV1() *metricRouterV1 {
	s := service.NewMetricServiceV1()

	gin.SetMode(gin.ReleaseMode)
	mux := gin.Default()
	mux.RedirectTrailingSlash = false
	mux.LoadHTMLFiles("web/template/allMetrics.html")

	mux.POST("/update/:metricType/:metricName/:metricValue", func(c *gin.Context) {
		c.Header("application-type", "text/plain")

		metricType := c.Param("metricType")
		metricName := c.Param("metricName")
		metricValue := c.Param("metricValue")

		switch metricType {
		case metricTypeGauge:
			metricValueFloat64, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
			}
			err = s.UpdateGauge(metricName, metricValueFloat64)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		case metricTypeCounter:
			metricValueInt64, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
			}
			s.UpdateCounter(metricName, metricValueInt64)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		default:
			c.AbortWithStatus(http.StatusNotImplemented)
		}
	})

	mux.GET("/value/:metricType/:metricName", func(c *gin.Context) {
		c.Header("application-type", "text/plain")

		metricType := c.Param("metricType")
		metricName := c.Param("metricName")

		switch metricType {
		case metricTypeGauge:
			metricValue, err := s.GetGauge(metricName)
			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
			}
			c.String(http.StatusOK, strconv.FormatFloat(metricValue, 'f', -1, 64))
		case metricTypeCounter:
			metricValue, err := s.GetCounter(metricName)
			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
			}
			c.String(http.StatusOK, strconv.FormatInt(metricValue, 10))
		default:
			c.AbortWithStatus(http.StatusNotFound)
		}
	})

	mux.GET("/", func(c *gin.Context) {
		c.Header("application-type", "text/plain")

		gaugeNameToValue, counterNameToValue, err := s.GetMetricAll()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		c.HTML(http.StatusOK, "allMetrics.html", gin.H{
			"GaugeMap":   gaugeNameToValue,
			"CounterMap": counterNameToValue,
		})
	})

	metricRouter := &metricRouterV1{
		mux:     mux,
		service: service.NewMetricServiceV1(),
	}

	return metricRouter
}
