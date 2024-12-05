package redis

import (
	"api-gateway/configs"
	"api-gateway/pkg/logger"
	"context"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
}

func NewRedisClient() *Redis {
	rd := connectRedis()
	return &Redis{rd}
}

func connectRedis() *redis.Client {
	rd := redis.NewClient(&redis.Options{
		Addr:     configs.C.Redis.Addr,
		Password: configs.C.Redis.Password,
		DB:       configs.C.Redis.DB, // 使用默認的資料庫
	})

	_, err := rd.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatal("無法連接到Redis: ", err)
		return nil
	}
	return rd

}

func (rd *Redis) CloseRedis() {
	defer rd.Client.Close()
}
