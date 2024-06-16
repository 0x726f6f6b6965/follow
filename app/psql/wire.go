//go:build wireinject
// +build wireinject

package psql

import (
	"github.com/0x726f6f6b6965/follow/internal/config"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitUserService(cfg *config.AppConfig) (service *gorm.DB, cleanup func(), err error) {
	panic(wire.Build(dbSet))
}
