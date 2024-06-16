package user

import (
	"github.com/0x726f6f6b6965/follow/internal/config"
	"github.com/0x726f6f6b6965/follow/internal/services"
	"github.com/0x726f6f6b6965/follow/internal/storage"
	"github.com/0x726f6f6b6965/follow/internal/storage/user"
	"github.com/0x726f6f6b6965/follow/pkg/logger"

	"github.com/google/wire"
)

var userService = wire.NewSet(userStorage, loggerSet, services.NewUserService)

var userStorage = wire.NewSet(dbSet, user.New)

var dbSet = wire.NewSet(dbConfig, storage.NewPostgres)

var loggerSet = wire.NewSet(logConfig, logger.NewLogger)

func dbConfig(cfg *config.AppConfig) *config.DBConfig {
	return &cfg.DB
}

func logConfig(cfg *config.AppConfig) *logger.LogConfig {
	return &cfg.Log
}
