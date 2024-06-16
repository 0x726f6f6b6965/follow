package config

import "github.com/0x726f6f6b6965/follow/pkg/logger"

type AppConfig struct {
	Env      string           `yaml:"env" mapstructure:"env" cobra-usage:"the application environment" cobra-default:"dev"`
	GrpcPort uint64           `yaml:"grpc-port" mapstructure:"grpc-port"`
	Log      logger.LogConfig `yaml:"log" mapstructure:"log"`
	DB       DBConfig         `yaml:"db" mapstructure:"db"`
}
