//go:build wireinject
// +build wireinject

package follow

import (
	"github.com/0x726f6f6b6965/follow/internal/config"
	pbFollow "github.com/0x726f6f6b6965/follow/protos/follow/v1"
	"github.com/google/wire"
	boom "github.com/tylertreat/BoomFilters"
)

func InitFollowService(cfg *config.AppConfig, filter *boom.CountingBloomFilter) (service pbFollow.FollowServiceServer, cleanup func(), err error) {
	panic(wire.Build(followService))
}
