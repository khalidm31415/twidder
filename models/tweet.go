package models

import (
	"time"

	"gorm.io/gorm"
)

type Tweet struct {
	ID               uint
	UserID           uint   `gorm:"not null"`
	User             User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Text             string `gorm:"type:varchar(280);not null;index:,class:FULLTEXT"`
	InReplyToTweetID *uint
	InReplyToTweet   *Tweet
	CreatedAt        time.Time `gorm:"not null"`
	UpdatedAt        time.Time `gorm:"not null"`
	DeletedAt        gorm.DeletedAt
}
