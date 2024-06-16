package api

import (
	"github.com/0x726f6f6b6965/follow/app/api/router"
	"github.com/0x726f6f6b6965/follow/app/api/services"
	"github.com/google/wire"
)

var routerSet = wire.NewSet(services.NewFollowAPI, services.NewUserAPI, router.NewRouter)
