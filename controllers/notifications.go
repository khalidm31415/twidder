package controllers

import (
	"net/http"
	"twidder/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func Notifications(c *gin.Context) {
	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)

	var notifications []models.Notification
	models.DB.Order("id desc").
		Scopes(Paginate(c)).
		Preload(clause.Associations).
		Preload("Follow.Follower").
		Preload("Follow.Followee").
		Preload("ReplyTweet.User").
		Preload("ReplyTweet.InReplyToTweet.User").
		Preload("Like.Tweet.User").
		Preload("Like.User").
		Where("user_id = ?", user.ID).Find(&notifications)
	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}
