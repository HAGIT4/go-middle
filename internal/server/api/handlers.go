package api

import (
	"net/http"
	"strconv"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/internal/server/storage"
	"github.com/HAGIT4/go-middle/pkg/models"
	"github.com/gin-gonic/gin"
)

const (
	plainTextParseMethodGet = iota
	plainTextParseMethodPost
	getResponseFormatJSON
	getResponseFormatPlain
)

func parseBatchJSONrequest() (h gin.HandlerFunc) {
	h = func(c *gin.Context) {
		contentHeader := c.Request.Header.Get("Content-Type")
		if contentHeader != "application/json" {
			c.AbortWithError(http.StatusBadRequest, newAPINoJSONHeaderError())
			return
		}
		reqMetricSlice := &[]models.Metrics{}
		if err := c.BindJSON(reqMetricSlice); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.Set("requestSlice", reqMetricSlice)
	}
	return
}

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
		c.Set("requestModel", reqMetricModel)
	}
	return
}

// GetHandler godoc
// @Summary Get metric information
// @Description Get metric information by ID
// @Accept json
// @Produce json
// @Param metricType path string true "gauge or counter"
// @Param metricName path string true "name of metric"
// @Success 200 {object} models.Metrics
// @Router /value/{metricType}/{metricName} [get]
func getHandler(sv service.MetricServiceInterface, getResponseFormat int) (h gin.HandlerFunc) {
	h = func(c *gin.Context) {
		var reqMetric interface{}
		var found bool
		if reqMetric, found = c.Get("requestModel"); !found {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		reqMetricModel := reqMetric.(*models.Metrics)
		respMetricModel, err := sv.GetMetric(reqMetricModel)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		if err := sv.ComputeHash(respMetricModel); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
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

// GetAllDataHTMLhandler godoc
// @Summary Get all metric data
// @Description Get all metric data in HTML format
// @Success 200 {string} HTML
// @Router / [get]
func getAllDataHTMLhandler(sv service.MetricServiceInterface) (h gin.HandlerFunc) {
	return func(c *gin.Context) {
		c.Header("application-type", "text/plain")
		gaugeNameToValue, counterNameToValue, err := sv.GetMetricAll()
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

// UpdateHandler godoc
// @Summary Update a metric
// @Accept json
// @Produce json
// @Router /update/ [post]
// @Success 200
func updateHandler(sv service.MetricServiceInterface) (h gin.HandlerFunc) {
	h = func(c *gin.Context) {
		var reqMetricModel interface{}
		var found bool
		if reqMetricModel, found = c.Get("requestModel"); !found {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if err := sv.UpdateMetric(reqMetricModel.(*models.Metrics)); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	return
}

// UpdateBatchHandler godoc
// @Summary Update metrics batch
// @Accept json
// @Router /updates/ [post]
// @Success 200
func updateBatchHandler(sv service.MetricServiceInterface) (h gin.HandlerFunc) {
	h = func(c *gin.Context) {
		var reqMetricSliceGet interface{}
		var reqMetricSlice *[]models.Metrics
		var found bool
		if reqMetricSliceGet, found = c.Get("requestSlice"); !found {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		reqMetricSlice = reqMetricSliceGet.(*[]models.Metrics)
		if err := sv.UpdateBatch(reqMetricSlice); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}
	return
}

// DatabasePingHandler godoc
// @Summary Ping database
// @Router /ping [get]
// @Success 200
func databasePingHandler(st storage.StorageInterface) (h gin.HandlerFunc) {
	h = func(c *gin.Context) {
		if err := st.Ping(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Status(http.StatusOK)
	}
	return
}
