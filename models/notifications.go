package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID           uint
	UserID       uint
	User         User
	Seen         bool
	FollowID     *uint
	Follow       *Follow
	ReplyTweetID *uint
	ReplyTweet   *Tweet
	LikeID       *uint
	Like         *Like
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
	DeletedAt    gorm.DeletedAt
}

func NewFollowedNotification(user User, follow Follow) *Notification {
	return &Notification{
		UserID:   user.ID,
		User:     user,
		Seen:     false,
		FollowID: &follow.ID,
		Follow:   &follow,
	}
}

func NewRepliedNotification(user User, replyTweet Tweet) *Notification {
	return &Notification{
		UserID:       user.ID,
		User:         user,
		Seen:         false,
		ReplyTweetID: &replyTweet.ID,
		ReplyTweet:   &replyTweet,
	}
}

func NewLikedNotification(user User, like Like) *Notification {
	return &Notification{
		UserID: user.ID,
		User:   user,
		Seen:   false,
		LikeID: &like.ID,
		Like:   &like,
	}
}
