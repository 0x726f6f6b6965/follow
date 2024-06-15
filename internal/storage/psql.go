package storage

import (
	"fmt"

	"github.com/0x726f6f6b6965/follow/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(cfg *config.DBConfig) (*gorm.DB, func(), error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName,
		cfg.Port, cfg.SSLmode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, func() {}, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, func() {}, err
	}
	return db, func() { sqlDB.Close() }, nil
}
