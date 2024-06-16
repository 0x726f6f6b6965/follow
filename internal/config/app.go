package config

import "github.com/0x726f6f6b6965/follow/pkg/logger"

type AppConfig struct {
	Env      string           `yaml:"env" mapstructure:"env" cobra-usage:"the application environment" cobra-default:"dev"`
	HttpPort uint64           `yaml:"http-port" mapstructure:"http-port"`
	Log      logger.LogConfig `yaml:"log" mapstructure:"log"`
	DB       DBConfig         `yaml:"db" mapstructure:"db"`
	Redis    RedisConfig      `yaml:"redis" mapstructure:"redis"`
	TTL      uint64           `yaml:"ttl" mapstructure:"ttl"`
}
