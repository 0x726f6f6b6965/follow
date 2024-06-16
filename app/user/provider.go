package user

import (
	"time"

	"github.com/0x726f6f6b6965/follow/internal/config"
	"github.com/0x726f6f6b6965/follow/internal/services"
	"github.com/0x726f6f6b6965/follow/internal/storage/cache"
	"github.com/0x726f6f6b6965/follow/internal/storage/user"
	"github.com/0x726f6f6b6965/follow/pkg/logger"

	"github.com/google/wire"
)

var userService = wire.NewSet(user.New, loggerSet, cache.New, getTTL, services.NewUserService)

var loggerSet = wire.NewSet(logConfig, logger.NewLogger)

func getTTL(cfg *config.AppConfig) time.Duration {
	return time.Duration(cfg.TTL) * time.Second
}

func logConfig(cfg *config.AppConfig) *logger.LogConfig {
	return &cfg.Log
}
