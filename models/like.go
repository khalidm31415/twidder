package models

import (
	"time"

	"gorm.io/gorm"
)

type Like struct {
	TweetID   uint      `gorm:"not null"`
	Tweet     Tweet     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt gorm.DeletedAt
}
