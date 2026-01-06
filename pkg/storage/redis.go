package storage

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/chatroom/utils"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

// GetRedis 获取Redis单例客户端
func GetRedis() *redis.Client {
	redisOnce.Do(func() {
		redisClient = initRedis()
	})
	return redisClient
}

func initRedis() *redis.Client {
	cfg, err := utils.LoadConfig()
	// 如果未配置 Redis，返回 nil（允许不使用 Redis 运行）
	if cfg.Redis.Addr == "" {
		log.Printf("Redis not configured, notification read status will use MySQL fallback")
		return nil
	}
	if err != nil {
		panic(err)
	}
	poolSize := cfg.Redis.PoolSize
	if poolSize <= 0 {
		poolSize = 10
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Addr, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: poolSize,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		log.Printf("Continuing without Redis, notification read status will use MySQL fallback")
		return nil
	}

	log.Printf("Redis connected: %s:%d", cfg.Redis.Addr, cfg.Redis.Port)
	return client
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
