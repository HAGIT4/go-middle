package middleware

import (
	"net/http"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/pkg/models"
	"github.com/gin-gonic/gin"
)

func CheckHashSHA256Middleware(sv service.MetricServiceInterface) (h gin.HandlerFunc) {
	h = func(c *gin.Context) {
		if len(sv.GetHashKey()) == 0 {
			return
		}
		var reqMetric interface{}
		reqMetric, found := c.Get("requestModel")
		if !found {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		reqMetricModel := reqMetric.(*models.Metrics)
		if err := sv.CheckHash(reqMetricModel); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}
	return
}

func CheckBatchHashSHA256Middleware(sv service.MetricServiceInterface) (h gin.HandlerFunc) {
	h = func(c *gin.Context) {
		if len(sv.GetHashKey()) == 0 {
			return
		}
		var reqMetricSliceGet interface{}
		var reqMetricSlice *[]models.Metrics
		var found bool
		if reqMetricSliceGet, found = c.Get("requestSlice"); !found {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		reqMetricSlice = reqMetricSliceGet.(*[]models.Metrics)
		for _, metric := range *reqMetricSlice {
			if err := sv.CheckHash(&metric); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
		}
	}
	return
}
