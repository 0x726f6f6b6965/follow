// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package rds

import (
	"github.com/0x726f6f6b6965/follow/internal/config"
	"github.com/0x726f6f6b6965/follow/internal/storage"
	"github.com/redis/go-redis/v9"
)

// Injectors from wire.go:

func InitRdsService(cfg *config.AppConfig) (*redis.Client, func(), error) {
	configRedisConfig := redisConfig(cfg)
	client, cleanup, err := storage.NewRedis(configRedisConfig)
	if err != nil {
		return nil, nil, err
	}
	return client, func() {
		cleanup()
	}, nil
}
