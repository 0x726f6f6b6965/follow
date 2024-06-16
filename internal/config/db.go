package config

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db-name"`
	SSLmode  string `yaml:"ssl-mode" default:"disable"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DBNum    int    `yaml:"db-num"`
	PoolSize int    `yaml:"pool-size"`
}
