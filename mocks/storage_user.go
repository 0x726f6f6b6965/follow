package mocks

import "github.com/0x726f6f6b6965/follow/internal/storage/models"

var (
	GetUserInfoFunc     func(usernames ...string) ([]models.User, error)
	GetUserInfoByIdFunc func(ids ...int) ([]models.User, error)
	CreateUserFunc      func(data *models.User) error
)

type MockSotrageUsers struct{}

func (m *MockSotrageUsers) GetUserInfo(usernames ...string) ([]models.User, error) {
	return GetUserInfoFunc(usernames...)
}
func (m *MockSotrageUsers) GetUserInfoById(ids ...int) ([]models.User, error) {
	return GetUserInfoByIdFunc(ids...)
}
func (m *MockSotrageUsers) CreateUser(data *models.User) error {
	return CreateUserFunc(data)
}
