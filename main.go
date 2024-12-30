package main

import (
	"api-gateway/configs"
	"api-gateway/grpc/notification"
	"api-gateway/grpc/user_auth"
	"api-gateway/grpc/user_info"
	"api-gateway/internal/endpoints"
	v1 "api-gateway/internal/endpoints/api/v1"
	"api-gateway/middleware"
	"api-gateway/pkg/logger"
	"api-gateway/pkg/redis"
	"api-gateway/routers"

	"go.uber.org/fx"
)

func init() {
	configs.Setup()
	logger.Setup()

}

func main() {

	app := fx.New(
		// 依賴注入
		fx.Provide(
			func() *redis.Redis {
				return redis.NewRedisClient(configs.C.Redis.Addr, configs.C.Redis.Password, configs.C.Redis.DB)
			},

			//middleware
			middleware.NewCheckJwtToken,
			middleware.NewRateLimiter,
			middleware.NewMiddleware,

			//grpc client
			func() user_auth.UserAuthServiceClient {
				return user_auth.NewUserAuthServiceClient(endpoints.CreateClient(configs.C.Service.UserAuth))
			},
			func() user_info.UserInfoServiceClient {
				return user_info.NewUserInfoServiceClient(endpoints.CreateClient(configs.C.Service.UserInfo))
			},
			func() notification.NotificationServiceClient {
				return notification.NewNotificationServiceClient(endpoints.CreateClient(configs.C.Service.Notification))
			},

			//api
			v1.NewUserInfo,
			v1.NewUserAuth,
			v1.NewNotify,

			//router
			routers.InitRouter,
		),

		// 啟動
		fx.Invoke(endpoints.StartServer),
	)

	app.Run()

}
