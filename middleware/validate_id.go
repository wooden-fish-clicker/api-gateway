package middleware

import (
	"api-gateway/pkg/app"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidateID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		appG := app.Gin{C: c}

		_, err := strconv.Atoi(id)
		if err != nil {
			appG.Response(http.StatusBadRequest, false, "path的id必須是number", "path的id必須是number", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
