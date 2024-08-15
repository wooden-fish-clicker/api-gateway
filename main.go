package main

import (
	"api-gateway/configs"
	"api-gateway/grpc/user_auth"
	"api-gateway/grpc/user_info"
	v1 "api-gateway/internal/api/v1"
	"api-gateway/pkg/logger"
	"api-gateway/pkg/redis"
	"api-gateway/routers"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	godotenv.Load(".env")
	configs.Setup()
	logger.Setup()
	// db.ConnectDB()
	redis.ConnectRedis()
	// db.ConnectMongoDB()
}

func main() {

	createGrpcClient(
		viper.GetString("USER_INFO_SERVICE_ADDR"),
		viper.GetString("USER_AUTH_SERVICE_ADDR"),
	)

	// Gin server
	gin.SetMode(configs.C.Server.RunMode)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", configs.C.Server.HttpPort),
		Handler:        routers.InitRouter(),
		ReadTimeout:    configs.C.Server.ReadTimeout * time.Second,
		WriteTimeout:   configs.C.Server.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	logger.Info("Starting server...")
	if err := s.ListenAndServe(); err != nil {
		logger.Fatal("Server failed to start: ", err)
	}

}

func createGrpcClient(microServices ...string) {

	for _, serviceAddr := range microServices {

		conn, err := grpc.NewClient(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal("無法連接到服務器 %s：%v", serviceAddr, err)
		}

		switch serviceAddr {
		case viper.GetString("USER_INFO_SERVICE_ADDR"):
			userInfoClient := user_info.NewUserInfoServiceClient(conn)
			v1.UserInfoService = userInfoClient
		case viper.GetString("USER_AUTH_SERVICE_ADDR"):
			userAuthClient := user_auth.NewUserAuthServiceClient(conn)
			v1.UserAuthService = userAuthClient
		}
	}
}
