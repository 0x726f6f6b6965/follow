//go:build wireinject
// +build wireinject

package redis

import (
	"github.com/0x726f6f6b6965/follow/internal/config"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

func InitUserService(cfg *config.AppConfig) (service *redis.Client, cleanup func(), err error) {
	panic(wire.Build(redisSet))
}
