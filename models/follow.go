package models

import (
	"time"

	"gorm.io/gorm"
)

type Follow struct {
	FollowerID uint      `gorm:"not null;uniqueIndex:idx_follow"`
	Follower   User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FolloweeID uint      `gorm:"not null;uniqueIndex:idx_follow"`
	Followee   User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
	DeletedAt  gorm.DeletedAt
}
