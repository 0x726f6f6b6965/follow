package user

import (
	"fmt"

	"github.com/0x726f6f6b6965/follow/internal/storage/models"
	"gorm.io/gorm"
)

type SotrageUsers struct {
	db *gorm.DB
}

func New(db *gorm.DB) *SotrageUsers {
	return &SotrageUsers{
		db: db,
	}
}

// // GetUserWithFollowers get user with followers
// func (m *SotrageUsers) GetUserWithFollowers(username string) (*models.UserFollowers, error) {
// 	data := &models.UserFollowers{}
// 	err := m.db.Preload("Followers").First(data, "username = ?", username).Error
// 	return data, err
// }

// // GetUserWithFollowing get user with following
// func (m *SotrageUsers) GetUserWithFollowing(username string) (*models.UserFollowing, error) {
// 	data := &models.UserFollowing{}
// 	err := m.db.Preload("Following").First(data, "username = ?", username).Error
// 	return data, err
// }

// // GetUserWithFollowers get user with followers
// func (m *SotrageUsers) GetUserWithFollowers(username string, lastId int, limit int) ([]models.Follower, error) {
// 	data := []models.Follower{}
// 	err := m.db.Table(fmt.Sprintf("%s f", models.FOLLOWERS_TABLE)).
// 		Joins(fmt.Sprintf("INNER JOIN %s u ON f.following_id = u.id", models.USER_TABLE)).
// 		Select("f.follower_id").
// 		Where("u.username = ? and f.follower_id > ?", username, lastId).Limit(limit).Find(&data).Error

// 	return data, err
// }

// GetUserWithFollowers get user with followers
func (m *SotrageUsers) GetUserWithFollowers(userId int, lastId int, limit int) ([]models.User, error) {
	data := []models.User{}
	err := m.db.Table(fmt.Sprintf("%s f", models.FOLLOWERS_TABLE)).
		Joins(fmt.Sprintf("INNER JOIN %s u ON f.follower_id = u.id", models.USER_TABLE)).
		Select("u.id, u.username").
		Where("f.following_id = ? and f.follower_id > ?", userId, lastId).
		Order("f.follower_id asc").
		Limit(limit).
		Find(&data).Error

	return data, err
}

// // GetUserWithFollowing get user with following
// func (m *SotrageUsers) GetUserWithFollowing(username string, lastId int, limit int) ([]models.Follower, error) {
// 	data := []models.Follower{}
// 	err := m.db.Table(fmt.Sprintf("%s f", models.FOLLOWERS_TABLE)).
// 		Joins(fmt.Sprintf("INNER JOIN %s u ON f.follower_id = u.id", models.USER_TABLE)).
// 		Select("f.following_id").
// 		Where("u.username = ? and f.following_id > ?", username, lastId).Limit(limit).Find(&data).Error

// 	return data, err
// }

// GetUserWithFollowing get user with following
func (m *SotrageUsers) GetUserWithFollowing(userId int, lastId int, limit int) ([]models.User, error) {
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

// // GetUserWithFriends get user with friends
// func (m *SotrageUsers) GetUserWithFriends(username string, lastId int, limit int) ([]models.Follower, error) {
// 	data := []models.Follower{}
// 	err := m.db.Table(fmt.Sprintf("%s f", models.FOLLOWERS_TABLE)).
// 		Joins(fmt.Sprintf("INNER JOIN %s f2 ON f.follower_id = f2.following_id and f.following_id = f2.follower_id", models.FOLLOWERS_TABLE)).
// 		Select("f.following_id").
// 		Where(
// 			fmt.Sprintf("f.follower_id < f.following_id and f.follower_id = (SELECT id FROM %s WHERE username = ?) and f.follower_id > ?", models.USER_TABLE),
// 			username, lastId).Limit(limit).Find(&data).Error

// 	return data, err
// }

// GetUserWithFriends get user with friends
func (m *SotrageUsers) GetUserWithFriends(userId int, lastId int, limit int) ([]models.Follower, error) {
	data := []models.Follower{}
	err := m.db.Table(fmt.Sprintf("%s f", models.FOLLOWERS_TABLE)).
		Joins(fmt.Sprintf("INNER JOIN %s f2 ON f.follower_id = f2.following_id and f.following_id = f2.follower_id", models.FOLLOWERS_TABLE)).
		Select("f.following_id").
		Where(
			"f.follower_id < f.following_id and f.follower_id = ? and f.follower_id > ?",
			userId, lastId).Limit(limit).Find(&data).Error

	return data, err
}

func (m *SotrageUsers) SetFollowing(userId int, targetId int) error {
	err := m.db.Exec(fmt.Sprintf("INSERT INTO %s (follower_id, following_id) VALUES (?, ?)", models.FOLLOWERS_TABLE), userId, targetId).Error
	return err
}

func (m *SotrageUsers) UnsetFollowing(userId int, targetId int) error {
	err := m.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE follower_id = ? AND following_id = ?", models.FOLLOWERS_TABLE), userId, targetId).Error
	return err
}

func (m *SotrageUsers) GetUserInfo(usernames ...string) ([]models.User, error) {
	data := []models.User{}
	err := m.db.Table(models.USER_TABLE).Select("id, username").Where("username IN ?", usernames).Find(&data).Error
	return data, err
}

func (m *SotrageUsers) GetUserInfoById(ids ...int) ([]models.User, error) {
	data := []models.User{}
	err := m.db.Table(models.USER_TABLE).Select("id, username").Where("id IN ?", ids).Find(&data).Error
	return data, err
}
