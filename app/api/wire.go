//go:build wireinject
// +build wireinject

package api

import (
	"github.com/0x726f6f6b6965/follow/app/api/router"
	pbFollow "github.com/0x726f6f6b6965/follow/protos/follow/v1"
	pbUser "github.com/0x726f6f6b6965/follow/protos/user/v1"
	"github.com/google/wire"
)

func InitRouter(user pbUser.UserServiceServer, follow pbFollow.FollowServiceServer) (services router.Router, err error) {
	panic(wire.Build(routerSet))
}
