package routers

import (
	v1 "api-gateway/internal/endpoints/api/v1"
	"api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RouterParams struct {
	fx.In

	Middleware *middleware.Middleware
	UserAuth   *v1.UserAuth
	UserInfo   *v1.UserInfo
	Notify     *v1.Notify
}

func InitRouter(routerParams RouterParams) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(routerParams.Middleware.RateLimiter.CheckRate())

	apiGroup := r.Group("/api/v1")

	apiGroup.POST("/user", routerParams.UserInfo.Register)
	apiGroup.POST("auth/login", routerParams.UserAuth.Login)

	apiGroup.Use(routerParams.Middleware.CheckJwtToken.CheckToken())
	apiGroup.PUT("/user", routerParams.UserInfo.UpdateUser)
	apiGroup.PATCH("/user", routerParams.UserInfo.UpdateUserPassword)
	apiGroup.GET("/user", routerParams.UserInfo.GetCurrentUserInfo)
	apiGroup.GET("/user/:id", routerParams.UserInfo.GetUserInfo)

	apiGroup.GET("auth/logout", routerParams.UserAuth.Logout)
	apiGroup.POST("auth/login/line", routerParams.UserAuth.LineLogin)

	apiGroup.GET("/notifications", routerParams.Notify.GetNotificationList)
	apiGroup.PUT("/notification/read", routerParams.Notify.ReadNotification)
	apiGroup.DELETE("/notification/:id", routerParams.Notify.DeleteNotification)

	return r
}
