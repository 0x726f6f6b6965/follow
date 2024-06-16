//go:build wireinject
// +build wireinject

package user

import (
	"github.com/0x726f6f6b6965/follow/internal/config"
	pbUser "github.com/0x726f6f6b6965/follow/protos/user/v1"
	"github.com/google/wire"
)

func InitUserService(cfg *config.AppConfig) (service pbUser.UserServiceServer, cleanup func(), err error) {
	panic(wire.Build(userService))
}
