//go:build wireinject
// +build wireinject

package follow

import (
	"github.com/0x726f6f6b6965/follow/internal/config"
	pbFollow "github.com/0x726f6f6b6965/follow/protos/follow/v1"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	boom "github.com/tylertreat/BoomFilters"
	"gorm.io/gorm"
)

func InitFollowService(cfg *config.AppConfig, db *gorm.DB, rds *redis.Client, filter *boom.CountingBloomFilter) (service pbFollow.FollowServiceServer, cleanup func(), err error) {
	panic(wire.Build(followService))
}
