package models

import "time"

type User struct {
	Id         int       `gorm:"column:id;type:integer;auto_increment;not null;primary_key" json:"id"`
	Username   string    `gorm:"column:username;type:varchar(128);not null" json:"username"`
	Password   string    `gorm:"column:password;type:varchar(128);not null" json:"password"`
	Salt       string    `gorm:"column:salt;type:varchar(64);not null" json:"salt"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"`
}

type UserFollowers struct {
	User
	Followers []*User `gorm:"many2many:t_followers;foreignKey:id;joinForeignKey:follower_id;JoinReferences:following_id;"`
}

type UserFollowing struct {
	User
	Following []*User `gorm:"many2many:t_followers;foreignKey:id;joinForeignKey:following_id;JoinReferences:follower_id"`
}
