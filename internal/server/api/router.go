package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/HAGIT4/go-middle/internal/server/models"
	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/gin-gonic/gin"
)

type metricRouter struct {
	mux     *gin.Engine
	service service.IMetricService
}

func newMetricRouter() *metricRouter {
	s := service.NewMetricService()

	mux := gin.Default()
	mux.RedirectTrailingSlash = false
	mux.POST("/update/:metricType/:metricName/:metricValue", func(c *gin.Context) {
		metricType := c.Param("metricType")
		metricName := c.Param("metricName")
		metricValue := c.Param("metricValue")

		switch metricType {
		case metricTypeGauge:
			metricValueFloat64, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
			}
			metricInfo := &models.MetricGaugeInfo{
				Name:  metricName,
				Value: metricValueFloat64,
			}
			s.UpdateGauge(metricInfo)
		case metricTypeCounter:
			metricValueInt64, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
			}
			metricInfo := &models.MetricCounterInfo{
				Name:  metricName,
				Value: metricValueInt64,
			}
			s.UpdateCounter(metricInfo)
		default:
			c.AbortWithStatus(http.StatusNotImplemented)
		}
	})

	mux.GET("/value/:metricType/:metricName", func(c *gin.Context) {
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
		gaugeNameToValue, counterNameToValue, err := s.GetMetricAll()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		templateData := &models.GetAllMetricsData{
			GaugeNameToValue:   gaugeNameToValue,
			CounterNameToValue: counterNameToValue,
		}
		fmt.Println(templateData)

		absPath, err := filepath.Abs("./web/template/allMetrics.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		pageTemplate, err := ioutil.ReadFile(absPath)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		page, err := template.New("All Metrics").Parse(string(pageTemplate))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		c.HTML(http.StatusOK, "allMetrics.html", page)
		fmt.Println(page)
	})

	metricRouter := &metricRouter{
		mux:     mux,
		service: service.NewMetricService(),
	}

	return metricRouter
}
