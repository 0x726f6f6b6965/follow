package redis

import (
	"time"

	"github.com/0x726f6f6b6965/follow/internal/config"
	"github.com/0x726f6f6b6965/follow/internal/storage"

	"github.com/google/wire"
)

var redisSet = wire.NewSet(redisConfig, storage.NewRedis)

func redisConfig(cfg *config.AppConfig) *config.RedisConfig {
	return &cfg.Redis
}

func getTTL(cfg *config.AppConfig) time.Duration {
	return time.Duration(cfg.TTL) * time.Second
}
