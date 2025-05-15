package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis(ctx context.Context) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ошибка подключения к Redis: %v", err)
	}

	fmt.Println("Redis подключён успешно")
	return nil
}

func GetRedisClient() *redis.Client {
	return redisClient
}
