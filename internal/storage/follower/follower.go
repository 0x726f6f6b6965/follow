package follower

import (
	"fmt"

	"github.com/0x726f6f6b6965/follow/internal/storage/models"
	"gorm.io/gorm"
)

type sotrageFollowers struct {
	db *gorm.DB
}

type SotrageFollowers interface {
	GetUserWithFollowers(userId int, lastId int, limit int) ([]models.User, error)
	GetUserWithFollowing(userId int, lastId int, limit int) ([]models.User, error)
	GetUserWithFriends(userId int, lastId int, limit int) ([]models.Follower, error)
	SetFollowing(userId int, targetId int) error
	UnsetFollowing(userId int, targetId int) error
}

func New(db *gorm.DB) SotrageFollowers {
	return &sotrageFollowers{
		db: db,
	}
}

// GetUserWithFollowers get user with followers
func (m *sotrageFollowers) GetUserWithFollowers(userId int, lastId int, limit int) ([]models.User, error) {
	data := []models.User{}
	err := m.db.Table(fmt.Sprintf("%s f", models.FOLLOWERS_TABLE)).
		Joins(fmt.Sprintf("INNER JOIN %s u ON f.follower_id = u.id", models.USER_TABLE)).
		Select("u.id, u.username").
		Where("f.follower_id > ? and f.following_id = ?", lastId, userId).
		Order("f.follower_id asc").
		Limit(limit).
		Find(&data).Error

	return data, err
}

// GetUserWithFollowing get user with following
func (m *sotrageFollowers) GetUserWithFollowing(userId int, lastId int, limit int) ([]models.User, error) {
	data := []models.User{}
	err := m.db.Table(fmt.Sprintf("%s f", models.FOLLOWERS_TABLE)).
		Joins(fmt.Sprintf("INNER JOIN %s u ON f.following_id = u.id", models.USER_TABLE)).
		Select("u.id, u.username").
		Where("f.follower_id = ? and f.following_id > ?", userId, lastId).
		Order("f.following_id asc").
		Limit(limit).
		Find(&data).Error

	return data, err
}

// GetUserWithFriends get user with friends
func (m *sotrageFollowers) GetUserWithFriends(userId int, lastId int, limit int) ([]models.Follower, error) {
	data := []models.Follower{}
	err := m.db.Table(fmt.Sprintf("%s f", models.FOLLOWERS_TABLE)).
		Joins(fmt.Sprintf("INNER JOIN %s f2 ON f.follower_id = f2.following_id and f.following_id = f2.follower_id", models.FOLLOWERS_TABLE)).
		Select("f.following_id").
		Where(
			"f.follower_id < f.following_id and f.follower_id = ? and f.following_id > ?",
			userId, lastId).
		Order("f.following_id asc").
		Limit(limit).Find(&data).Error

	return data, err
}

func (m *sotrageFollowers) SetFollowing(userId int, targetId int) error {
	err := m.db.Exec(fmt.Sprintf("INSERT INTO %s (follower_id, following_id) VALUES (?, ?)", models.FOLLOWERS_TABLE), userId, targetId).Error
	return err
}

func (m *sotrageFollowers) UnsetFollowing(userId int, targetId int) error {
	err := m.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE follower_id = ? AND following_id = ?", models.FOLLOWERS_TABLE), userId, targetId).Error
	return err
}
