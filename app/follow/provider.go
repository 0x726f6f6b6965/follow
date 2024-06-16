package follow

import (
	"time"

	"github.com/0x726f6f6b6965/follow/internal/config"
	"github.com/0x726f6f6b6965/follow/internal/services"
	"github.com/0x726f6f6b6965/follow/internal/storage/cache"
	"github.com/0x726f6f6b6965/follow/internal/storage/follower"
	"github.com/0x726f6f6b6965/follow/internal/storage/user"
	"github.com/0x726f6f6b6965/follow/pkg/logger"

	"github.com/google/wire"
)

var followService = wire.NewSet(user.New, follower.New, cache.New, loggerSet, getTTL, services.NewFollowService)

var loggerSet = wire.NewSet(logConfig, logger.NewLogger)

func logConfig(cfg *config.AppConfig) *logger.LogConfig {
	return &cfg.Log
}

func getTTL(cfg *config.AppConfig) time.Duration {
	return time.Duration(cfg.TTL) * time.Second
}
