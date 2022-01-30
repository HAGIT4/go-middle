package middleware

import (
	"net/http"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/pkg/models"
	"github.com/gin-gonic/gin"
)

func CheckHashSHA256Middleware(s service.MetricServiceInterface) (h gin.HandlerFunc) {
	h = func(c *gin.Context) {
		if len(s.GetHashKey()) == 0 {
			return
		}
		var reqMetric interface{}
		reqMetric, found := c.Get("requestModel")
		if !found {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		reqMetricModel := reqMetric.(*models.Metrics)
		if err := s.CheckHash(reqMetricModel); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
	}
	return
}
