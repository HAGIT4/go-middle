package middleware

import (
	"net"
	"net/http"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/gin-gonic/gin"
)

func CheckTrustedSubnetMiddleWare(sv service.MetricServiceInterface) (h gin.HandlerFunc) {
	h = func(c *gin.Context) {
		trustedSubnet := sv.GetTrustedSubnet()
		if trustedSubnet == nil {
			return
		}
		realIP := net.ParseIP(c.GetHeader("X-Real-IP"))
		if realIP == nil {
			c.AbortWithStatus(http.StatusForbidden)
		}
		if trustedSubnet.Contains(realIP) {
			return
		} else {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
	return
}
