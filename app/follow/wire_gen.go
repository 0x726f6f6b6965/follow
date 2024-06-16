// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package follow

import (
	"github.com/0x726f6f6b6965/follow/internal/config"
	"github.com/0x726f6f6b6965/follow/internal/services"
	"github.com/0x726f6f6b6965/follow/internal/storage"
	"github.com/0x726f6f6b6965/follow/internal/storage/user"
	"github.com/0x726f6f6b6965/follow/pkg/logger"
	"github.com/0x726f6f6b6965/follow/protos/follow/v1"
)

// Injectors from wire.go:

func InitFollowService(cfg *config.AppConfig) (v1.FollowServiceServer, func(), error) {
	configDBConfig := dbConfig(cfg)
	db, cleanup, err := storage.NewPostgres(configDBConfig)
	if err != nil {
		return nil, nil, err
	}
	sotrageUsers := user.New(db)
	sotrageFollowers, cleanup2, err := getFollowStorage(cfg)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	loggerLogConfig := logConfig(cfg)
	zapLogger, cleanup3, err := logger.NewLogger(loggerLogConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	followServiceServer := services.NewFollowService(sotrageUsers, sotrageFollowers, zapLogger)
	return followServiceServer, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
