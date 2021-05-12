package controllers

import (
	"errors"
	"fmt"
	"gin-twitter/models"
	"net/http"

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

	var input LikeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var tweet models.Tweet
	result := models.DB.First(&tweet, input.TweetID)
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
	c.JSON(http.StatusCreated, gin.H{"payload": result})
}

func Unlike(c *gin.Context) {
	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)

	var input UnlikeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	result := models.DB.Delete(models.Like{}, "user_id = ? AND tweet_id = ?", user.ID, input.TweetID)
	if result.Error != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"payload": result})
}
