package main

import (
	"api-gateway/configs"
	"api-gateway/grpc/notification"
	"api-gateway/grpc/user_auth"
	"api-gateway/grpc/user_info"
	v1 "api-gateway/internal/api/v1"
	"api-gateway/middleware"
	"api-gateway/pkg/logger"
	"api-gateway/pkg/redis"
	"api-gateway/routers"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	configs.Setup()
	logger.Setup()

}

func main() {

	app := fx.New(
		// 依賴注入
		fx.Provide(
			redis.NewRedisClient,

			//middleware
			middleware.NewCheckJwtToken,
			middleware.NewRateLimiter,
			middleware.NewMiddleware,

			//grpc client
			func() user_auth.UserAuthServiceClient {
				return user_auth.NewUserAuthServiceClient(createClient(configs.C.Service.UserAuth))
			},
			func() user_info.UserInfoServiceClient {
				return user_info.NewUserInfoServiceClient(createClient(configs.C.Service.UserInfo))
			},
			func() notification.NotificationServiceClient {
				return notification.NewNotificationServiceClient(createClient(configs.C.Service.Notification))
			},

			//api
			v1.NewUserInfo,
			v1.NewUserAuth,
			v1.NewNotify,

			//router
			routers.InitRouter,
		),

		// 啟動
		fx.Invoke(startServer),
	)

	app.Run()

}

func createClient(serviceAddr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("無法連接到服務器 %s:%v", serviceAddr, err)
	}
	return conn
}

func startServer(lc fx.Lifecycle, router *gin.Engine) {
	gin.SetMode(configs.C.Server.RunMode)
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", configs.C.Server.HttpPort),
		Handler:        router,
		ReadTimeout:    configs.C.Server.ReadTimeout * time.Second,
		WriteTimeout:   configs.C.Server.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting server on :8080")

			go func() {
				if err := server.ListenAndServe(); err != nil {
					logger.Fatal("Starting server error ：%v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping server")
			return server.Shutdown(ctx)
		},
	})
}
