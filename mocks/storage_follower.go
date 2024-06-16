package mocks

import "github.com/0x726f6f6b6965/follow/internal/storage/models"

var (
	GetUserWithFollowersFunc func(userId int, lastId int, limit int) ([]models.User, error)
	GetUserWithFollowingFunc func(userId int, lastId int, limit int) ([]models.User, error)
	GetUserWithFriendsFunc   func(userId int, lastId int, limit int) ([]models.Follower, error)
	SetFollowingFunc         func(userId int, targetId int) error
	UnsetFollowingFunc       func(userId int, targetId int) error
)

type MockSotrageFollowers struct{}

func (m *MockSotrageFollowers) GetUserWithFollowers(userId int, lastId int, limit int) ([]models.User, error) {
	return GetUserWithFollowersFunc(userId, lastId, limit)
}
func (m *MockSotrageFollowers) GetUserWithFollowing(userId int, lastId int, limit int) ([]models.User, error) {
	return GetUserWithFollowingFunc(userId, lastId, limit)
}
func (m *MockSotrageFollowers) GetUserWithFriends(userId int, lastId int, limit int) ([]models.Follower, error) {
	return GetUserWithFriendsFunc(userId, lastId, limit)
}
func (m *MockSotrageFollowers) SetFollowing(userId int, targetId int) error {
	return SetFollowingFunc(userId, targetId)
}
func (m *MockSotrageFollowers) UnsetFollowing(userId int, targetId int) error {
	return UnsetFollowingFunc(userId, targetId)
}
