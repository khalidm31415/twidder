package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"twidder/models"

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
	c.JSON(http.StatusCreated, gin.H{"liked": tweet})
}

func Unlike(c *gin.Context) {
	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)
	tweetId, _ := strconv.Atoi(c.Param("id"))

	var tweet models.Tweet
	queryResult := models.DB.Preload("User").Take(&tweet, tweetId)
	if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
		fmt.Println(fmt.Errorf("[ERROR] %v", queryResult.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	deleteResult := models.DB.Delete(models.Like{}, "user_id = ? AND tweet_id = ?", user.ID, tweetId)
	if deleteResult.Error != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", deleteResult.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"unliked": tweet})
}

func Likes(c *gin.Context) {
	tweetId, _ := strconv.Atoi(c.Param("id"))

	var likes []models.Like
	result := models.DB.Preload("User").Where("tweet_id = ?", tweetId).Find(&likes)
	if result.Error != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	users := []models.User{}
	for _, like := range likes {
		users = append(users, like.User)
	}

	c.JSON(http.StatusOK, gin.H{"likes": users})
}
