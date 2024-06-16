//go:build wireinject
// +build wireinject

package user

import (
	"github.com/0x726f6f6b6965/follow/internal/config"
	pbUser "github.com/0x726f6f6b6965/follow/protos/user/v1"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	boom "github.com/tylertreat/BoomFilters"
	"gorm.io/gorm"
)

func InitUserService(cfg *config.AppConfig, db *gorm.DB, rds *redis.Client, filter *boom.CountingBloomFilter) (service pbUser.UserServiceServer, cleanup func(), err error) {
	panic(wire.Build(userService))
}
