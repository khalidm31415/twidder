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

type CreateTweetInput struct {
	Text string `binding:"required"`
}

type ReplyInput struct {
	Text string `binding:"required"`
}

func FindTweets(c *gin.Context) {
	tx := models.DB.Preload("User")

	q := c.Query("q")
	if len(q) > 0 {
		tx.Where("MATCH(text) AGAINST(? IN NATURAL LANGUAGE MODE)", q)
	}

	var tweets []models.Tweet
	tx.Scopes(Paginate(c)).Find(&tweets)
	c.JSON(http.StatusOK, gin.H{"tweets": tweets})
}

func FindTweet(c *gin.Context) {
	tweetId, _ := strconv.Atoi(c.Param("id"))
	var tweet models.Tweet
	result := models.DB.Preload("User").Take(&tweet, tweetId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println(fmt.Errorf("[ERROR] %v", result.Error))
		c.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
		return
	}

	conversation := []models.Tweet{}
	currentTweetId := tweet.InReplyToTweetID
	for currentTweetId != nil {
		var inReplyToTweet models.Tweet
		models.DB.Preload("User").Take(&inReplyToTweet, currentTweetId)
		conversation = append(conversation, inReplyToTweet)
		currentTweetId = inReplyToTweet.InReplyToTweetID
	}

	c.JSON(http.StatusOK, gin.H{"tweet": tweet, "conversation": conversation})
}

func Timeline(c *gin.Context) {
	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)

	follows := []models.Follow{}
	models.DB.Preload("Followee").Where("follower_id = ?", user.ID).Find(&follows)

	followingIds := []int{}
	for _, follow := range follows {
		followingIds = append(followingIds, int(follow.FolloweeID))
	}

	var tweets []models.Tweet
	models.DB.Preload("User").Where("user_id IN ?", followingIds).Scopes(Paginate(c)).Find(&tweets)
	c.JSON(http.StatusOK, gin.H{"tweets": tweets})
}

func CreateTweet(c *gin.Context) {
	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)

	var input CreateTweetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweet := models.Tweet{
		Text:   input.Text,
		UserID: user.ID,
		User:   user,
	}

	if result := models.DB.Create(&tweet); result.Error != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tweets": tweet})
}

func Reply(c *gin.Context) {
	inReplyToTweetId, _ := strconv.Atoi(c.Param("id"))
	var inReplyToTweet models.Tweet
	searchResult := models.DB.Preload("User").Take(&inReplyToTweet, inReplyToTweetId)
	if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
		fmt.Println(fmt.Errorf("[ERROR] %v", searchResult.Error))
		c.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
		return
	}

	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)

	var input CreateTweetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweet := models.Tweet{
		Text:             input.Text,
		UserID:           user.ID,
		User:             user,
		InReplyToTweetID: &inReplyToTweet.ID,
	}

	transactionErr := models.DB.Transaction(func(tx *gorm.DB) error {

		if createTweetResult := tx.Create(&tweet); createTweetResult.Error != nil {
			fmt.Println(fmt.Errorf("[ERROR] %v", createTweetResult.Error))
			return createTweetResult.Error
		}

		repliedNotifications := models.NewRepliedNotification(inReplyToTweet.User, tweet)
		if createNotificationResult := tx.Create(&repliedNotifications); createNotificationResult.Error != nil {
			fmt.Println(fmt.Errorf("[ERROR] %v", createNotificationResult.Error))
			return createNotificationResult.Error
		}

		return nil
	})

	if transactionErr != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", transactionErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"replied": tweet})
}

func Replies(c *gin.Context) {
	inReplyToTweetId, _ := strconv.Atoi(c.Param("id"))

	var replies []models.Tweet
	models.DB.Scopes(Paginate(c)).Where("in_reply_to_tweet_id = ?", inReplyToTweetId).Preload("User").Find(&replies)
	c.JSON(http.StatusOK, gin.H{"replies": replies})
}

func DeleteTweet(c *gin.Context) {
	tweetId, _ := strconv.Atoi(c.Param("id"))
	var tweet models.Tweet
	searchResult := models.DB.Take(&tweet, tweetId)
	if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
		fmt.Println(fmt.Errorf("[ERROR] %v", searchResult.Error))
		c.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
		return
	}

	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)
	if user.ID != uint(tweet.UserID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User cannot delete other user's tweet"})
		return
	}

	deleteResult := models.DB.Delete(&tweet)
	if deleteResult.Error != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", deleteResult.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"deleted": deleteResult})
}
