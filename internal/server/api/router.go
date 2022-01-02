package api

import (
	"net/http"
	"strconv"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/pkg/models"
	"github.com/gin-gonic/gin"
)

type metricRouter struct {
	mux     *gin.Engine
	service service.MetricServiceInterface
}

func newMetricRouter() *metricRouter {
	s := service.NewMetricService()

	mux := gin.Default()
	mux.RedirectTrailingSlash = false
	mux.LoadHTMLFiles("web/template/allMetrics.html")

	mux.POST("/update/", func(c *gin.Context) {
		contentHeader := c.Request.Header.Get("Content-Type")
		if contentHeader != "application/json" {
			c.AbortWithError(http.StatusBadRequest, newApiNoJSONHeaderError())
			return
		}
		reqMetricMsg := &models.Metrics{}
		if err := c.BindJSON(reqMetricMsg); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err := s.UpdateMetric(reqMetricMsg); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	})

	mux.POST("/update/:metricType/:metricName/:metricValue", func(c *gin.Context) {
		c.Header("application-type", "text/plain")
		metricType := c.Param("metricType")
		metricName := c.Param("metricName")
		metricValue := c.Param("metricValue")

		switch metricType {
		case metricTypeGauge:
			metricValueFloat64, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			err = s.UpdateGauge(metricName, metricValueFloat64)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		case metricTypeCounter:
			metricValueInt64, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			s.UpdateCounter(metricName, metricValueInt64)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		default:
			c.AbortWithStatus(http.StatusNotImplemented)
			return
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
				c.AbortWithError(http.StatusNotFound, err)
				return
			}
			c.String(http.StatusOK, strconv.FormatFloat(metricValue, 'f', -1, 64))
			return
		case metricTypeCounter:
			metricValue, err := s.GetCounter(metricName)
			if err != nil {
				c.AbortWithError(http.StatusNotFound, err)
				return
			}
			c.String(http.StatusOK, strconv.FormatInt(metricValue, 10))
			return
		default:
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	})

	mux.POST("/value/", func(c *gin.Context) {
		reqContentHeader := c.Request.Header.Get("Content-Type")
		if reqContentHeader != "application/json" {
			c.AbortWithError(http.StatusBadRequest, newApiNoJSONHeaderError())
			return
		}
		reqMetricMsg := &models.Metrics{}
		if err := c.BindJSON(reqMetricMsg); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		respMetricMsg, err := s.GetMetric(reqMetricMsg)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, respMetricMsg)
	})

	mux.GET("/", func(c *gin.Context) {
		c.Header("application-type", "text/plain")

		gaugeNameToValue, counterNameToValue, err := s.GetMetricAll()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.HTML(http.StatusOK, "allMetrics.html", gin.H{
			"GaugeMap":   gaugeNameToValue,
			"CounterMap": counterNameToValue,
		})
	})

	metricRouter := &metricRouter{
		mux:     mux,
		service: service.NewMetricService(),
	}

	return metricRouter
}
