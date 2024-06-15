package config

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db-name"`
	SSLmode  string `yaml:"ssl-mode" default:"disable"`
}
