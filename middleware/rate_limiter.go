package middleware

import (
	"api-gateway/pkg/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 每秒20，最多累積40個
var limiter = rate.NewLimiter(20, 40)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		if !limiter.Allow() {
			appG.Response(http.StatusTooManyRequests, false, "請求過多", nil, nil)
		}

		c.Next()
	}
}
