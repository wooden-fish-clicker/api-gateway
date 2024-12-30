package endpoints

import (
	"api-gateway/configs"
	"api-gateway/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartServer(lc fx.Lifecycle, router *gin.Engine) {
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

func CreateClient(serviceAddr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("無法連接到服務器 %s:%v", serviceAddr, err)
	}
	return conn
}
