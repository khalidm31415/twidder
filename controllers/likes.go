package controllers

import (
	"errors"
	"fmt"
	"gin-twitter/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LikeInput struct {
	TweetID uint `json:"tweet_id" binding:"required"`
}

type UnlikeInput struct {
	TweetID uint `json:"tweet_id" binding:"required"`
}

func Like(c *gin.Context) {
	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)
	tweetId, _ := strconv.Atoi(c.Param("id"))

	var tweet models.Tweet
	result := models.DB.First(&tweet, tweetId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println(fmt.Errorf("[ERROR] %v", result.Error))
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error})
		return
	}

	like := models.Like{
		UserID:  user.ID,
		User:    user,
		TweetID: tweet.ID,
		Tweet:   tweet,
	}

	if result := models.DB.Create(&like); result.Error != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"liked": result})
}

func Unlike(c *gin.Context) {
	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)
	tweetId, _ := strconv.Atoi(c.Param("id"))

	result := models.DB.Delete(models.Like{}, "user_id = ? AND tweet_id = ?", user.ID, tweetId)
	if result.Error != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"unliked": result})
}
