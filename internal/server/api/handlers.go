package api

import (
	"net/http"
	"strconv"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/pkg/models"
	"github.com/gin-gonic/gin"
)

const (
	plainTextParseMethodGet = iota
	plainTextParseMethodPost
	getResponseFormatJSON
	getResponseFormatPlain
)

func parseJSONrequest() (h gin.HandlerFunc) {
	h = func(c *gin.Context) {
		contentHeader := c.Request.Header.Get("Content-Type")
		if contentHeader != "application/json" {
			c.AbortWithError(http.StatusBadRequest, newAPINoJSONHeaderError())
			return
		}
		reqMetricModel := &models.Metrics{}
		if err := c.BindJSON(reqMetricModel); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.Set("requestModel", reqMetricModel)
	}
	return
}

func parsePlainTextRequest(parseMethod int) (h gin.HandlerFunc) {
	h = func(c *gin.Context) {
		metricType := c.Param("metricType")
		metricName := c.Param("metricName")
		reqMetricModel := &models.Metrics{
			ID:    metricName,
			MType: metricType,
		}

		if parseMethod == plainTextParseMethodPost {
			metricValue := c.Param("metricValue")
			switch metricType {
			case metricTypeGauge:
				metricValueFloat64, err := strconv.ParseFloat(metricValue, 64)
				if err != nil {
					c.AbortWithError(http.StatusBadRequest, err)
					return
				}
				reqMetricModel.Value = &metricValueFloat64
			case metricTypeCounter:
				metricValueInt64, err := strconv.ParseInt(metricValue, 10, 64)
				if err != nil {
					c.AbortWithError(http.StatusBadRequest, err)
					return
				}
				reqMetricModel.Delta = &metricValueInt64
			default:
				c.AbortWithStatus(http.StatusNotImplemented)
				return
			}
		}
		c.Header("application-type", "text/plain")
		c.Set("requestModel", reqMetricModel)
	}
	return
}

func getHandler(s service.MetricServiceInterface, getResponseFormat int) (h gin.HandlerFunc) {
	h = func(c *gin.Context) {
		var reqMetric interface{}
		var found bool
		if reqMetric, found = c.Get("requestModel"); !found {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		reqMetricModel := reqMetric.(*models.Metrics)
		respMetricModel, err := s.GetMetric(reqMetricModel)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		if getResponseFormat == getResponseFormatJSON {
			c.JSON(http.StatusOK, *respMetricModel)
			return
		} else if getResponseFormat == getResponseFormatPlain {
			switch respMetricModel.MType {
			case metricTypeGauge:
				c.String(http.StatusOK, strconv.FormatFloat(*respMetricModel.Value, 'f', -1, 64))
				return
			case metricTypeCounter:
				c.String(http.StatusOK, strconv.FormatInt(*respMetricModel.Delta, 10))
				return
			default:
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		} else {
			return
		}
	}
	return
}

func getAllDataHTMLhandler(s service.MetricServiceInterface) (h gin.HandlerFunc) {
	return func(c *gin.Context) {
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
	}
}

func updateHandler(s service.MetricServiceInterface) (h gin.HandlerFunc) {
	h = func(c *gin.Context) {
		var reqMetricModel interface{}
		var found bool
		if reqMetricModel, found = c.Get("requestModel"); !found {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		if err := s.UpdateMetric(reqMetricModel.(*models.Metrics)); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
	return
}

func updateByJSONhandler(s service.MetricServiceInterface) (h gin.HandlerFunc) {
	return func(c *gin.Context) {
		contentHeader := c.Request.Header.Get("Content-Type")
		if contentHeader != "application/json" {
			c.AbortWithError(http.StatusBadRequest, newAPINoJSONHeaderError())
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
	}
}

func updateByResolveHandler(s service.MetricServiceInterface) (h gin.HandlerFunc) {
	return func(c *gin.Context) {
		c.Header("application-type", "text/plain")
		metricType := c.Param("metricType")
		metricName := c.Param("metricName")
		metricValue := c.Param("metricValue")

		var reqMetricMsg *models.Metrics
		switch metricType {
		case metricTypeGauge:
			metricValueFloat64, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			reqMetricMsg = &models.Metrics{
				ID:    metricName,
				MType: metricTypeGauge,
				Value: &metricValueFloat64,
			}
		case metricTypeCounter:
			metricValueInt64, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			reqMetricMsg = &models.Metrics{
				ID:    metricName,
				MType: metricTypeCounter,
				Delta: &metricValueInt64,
			}
		default:
			c.AbortWithStatus(http.StatusNotImplemented)
			return
		}

		if err := s.UpdateMetric(reqMetricMsg); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}
}

func getByResolveHandler(s service.MetricServiceInterface) (h gin.HandlerFunc) {
	return func(c *gin.Context) {
		c.Header("application-type", "text/plain")

		metricType := c.Param("metricType")
		metricName := c.Param("metricName")

		reqMetricMsg := &models.Metrics{
			ID:    metricName,
			MType: metricType,
		}

		respMetricMsg, err := s.GetMetric(reqMetricMsg)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		switch respMetricMsg.MType {
		case metricTypeGauge:
			c.String(http.StatusOK, strconv.FormatFloat(*respMetricMsg.Value, 'f', -1, 64))
			return
		case metricTypeCounter:
			c.String(http.StatusOK, strconv.FormatInt(*respMetricMsg.Delta, 10))
			return
		default:
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}
}

func getByJSONhandler(s service.MetricServiceInterface) (h gin.HandlerFunc) {
	return func(c *gin.Context) {
		reqContentHeader := c.Request.Header.Get("Content-Type")
		if reqContentHeader != "application/json" {
			c.AbortWithError(http.StatusBadRequest, newAPINoJSONHeaderError())
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
		c.JSON(http.StatusOK, respMetricMsg)
	}
}
