package models

import "time"

type Follower struct {
	Id          int       `gorm:"column:id;type:integer;primary_key" json:"id"`
	FollowerId  int       `gorm:"column:follower_id;type:integer;not null" json:"follower_id"`
	FollowingId int       `gorm:"column:following_id;type:integer;not null" json:"following_id"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"`
}
