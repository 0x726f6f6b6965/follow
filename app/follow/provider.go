package follow

import (
	"time"

	"github.com/0x726f6f6b6965/follow/internal/config"
	"github.com/0x726f6f6b6965/follow/internal/services"
	"github.com/0x726f6f6b6965/follow/internal/storage"
	"github.com/0x726f6f6b6965/follow/internal/storage/cache"
	"github.com/0x726f6f6b6965/follow/internal/storage/follower"
	"github.com/0x726f6f6b6965/follow/internal/storage/user"
	"github.com/0x726f6f6b6965/follow/pkg/logger"

	"github.com/google/wire"
)

var followService = wire.NewSet(userStorage, getFollowStorage, loggerSet, redisSet, getTTL, services.NewFollowService)

var userStorage = wire.NewSet(dbSet, user.New)

var dbSet = wire.NewSet(dbConfig, storage.NewPostgres)

var redisSet = wire.NewSet(redisConfig, storage.NewRedis, cache.New)

var loggerSet = wire.NewSet(logConfig, logger.NewLogger)

func dbConfig(cfg *config.AppConfig) *config.DBConfig {
	return &cfg.DB
}

func logConfig(cfg *config.AppConfig) *logger.LogConfig {
	return &cfg.Log
}

func redisConfig(cfg *config.AppConfig) *config.RedisConfig {
	return &cfg.Redis
}

func getTTL(cfg *config.AppConfig) time.Duration {
	return time.Duration(cfg.TTL) * time.Second
}

func getFollowStorage(cfg *config.AppConfig) (follower.SotrageFollowers, func(), error) {
	db, cleanup, err := storage.NewPostgres(&cfg.DB)
	return follower.New(db), cleanup, err
}
