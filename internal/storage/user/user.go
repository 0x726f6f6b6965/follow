package user

import (
	"github.com/0x726f6f6b6965/follow/internal/storage/models"
	"gorm.io/gorm"
)

type sotrageUsers struct {
	db *gorm.DB
}

type SotrageUsers interface {
	GetUserInfo(usernames ...string) ([]models.User, error)
	GetUserInfoById(ids ...int) ([]models.User, error)
	CreateUser(data *models.User) error
}

func New(db *gorm.DB) SotrageUsers {
	return &sotrageUsers{
		db: db,
	}
}

func (m *sotrageUsers) GetUserInfo(usernames ...string) ([]models.User, error) {
	data := []models.User{}
	err := m.db.Table(models.USER_TABLE).Select("id, username").Where("username IN ?", usernames).Find(&data).Error
	return data, err
}

func (m *sotrageUsers) GetUserInfoById(ids ...int) ([]models.User, error) {
	data := []models.User{}
	err := m.db.Table(models.USER_TABLE).Select("id, username").Where("id IN ?", ids).Find(&data).Error
	return data, err
}

func (m *sotrageUsers) CreateUser(data *models.User) error {
	return m.db.Table(models.USER_TABLE).Create(data).Error
}
