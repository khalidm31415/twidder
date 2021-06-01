package models

import (
	"time"

	"gorm.io/gorm"
)

type Like struct {
	ID        uint
	TweetID   uint      `gorm:"not null;uniqueIndex:idx_like"`
	Tweet     Tweet     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_like"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"not null"`
	DeletedAt gorm.DeletedAt
}
