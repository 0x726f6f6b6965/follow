package psql

import (
	"github.com/0x726f6f6b6965/follow/internal/config"
	"github.com/0x726f6f6b6965/follow/internal/storage"

	"github.com/google/wire"
)

var dbSet = wire.NewSet(dbConfig, storage.NewPostgres)

func dbConfig(cfg *config.AppConfig) *config.DBConfig {
	return &cfg.DB
}
