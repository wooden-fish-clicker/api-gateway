package routers

import (
	apiV1 "api-gateway/internal/api/v1"
	"api-gateway/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RateLimiter())

	apiGroup := r.Group("/api/v1")

	apiGroup.POST("/user", apiV1.Register)
	apiGroup.POST("auth/login", apiV1.Login)

	apiGroup.Use(middleware.CheckToken())
	apiGroup.PUT("/user", apiV1.UpdateUser)
	apiGroup.PATCH("/user", apiV1.UpdateUserPassword)
	apiGroup.GET("/user", apiV1.GetCurrentUserInfo)
	apiGroup.GET("/user/:id", apiV1.GetUserInfo)

	apiGroup.GET("auth/logout", apiV1.Logout)
	apiGroup.POST("auth/login/line", apiV1.LineLogin)

	apiGroup.GET("/notifications", apiV1.GetNotificationList)
	apiGroup.PUT("/notification/read", apiV1.ReadNotification)
	apiGroup.DELETE("/notification/:id", apiV1.DeleteNotification)

	return r
}
