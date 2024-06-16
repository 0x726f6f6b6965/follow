//go:build wireinject
// +build wireinject

package follow

import (
	"github.com/0x726f6f6b6965/follow/internal/config"
	pbFollow "github.com/0x726f6f6b6965/follow/protos/follow/v1"
	"github.com/google/wire"
)

func InitFollowService(cfg *config.AppConfig) (service pbFollow.FollowServiceServer, cleanup func(), err error) {
	panic(wire.Build(followService))
}
