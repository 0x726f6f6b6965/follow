package storage

import (
	"fmt"

	"github.com/0x726f6f6b6965/follow/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.RedisConfig) (*redis.Client, func(), error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DBNum,
		PoolSize: cfg.PoolSize,
	})
	return redisClient, func() { redisClient.Close() }, nil
}
