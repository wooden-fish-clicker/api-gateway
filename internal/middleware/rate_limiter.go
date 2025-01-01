package middleware

import (
	"api-gateway/pkg/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	Limiter *rate.Limiter
}

func NewRateLimiter() *RateLimiter {
	// 每秒最多允許20個請求，最多累積40個
	limiter := rate.NewLimiter(20, 40)
	return &RateLimiter{Limiter: limiter}
}

func (rl *RateLimiter) CheckRate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		if !rl.Limiter.Allow() {
			appG.Response(http.StatusTooManyRequests, false, "請求過多", nil, nil)
		}

		c.Next()
	}
}
