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

func Follow(c *gin.Context) {
	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)

	followeeId, _ := strconv.Atoi(c.Param("id"))

	if user.ID == uint(followeeId) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Users cannot follow themselves"})
		return
	}

	var followee models.User
	result := models.DB.First(&followee, followeeId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println(fmt.Errorf("[ERROR]: %v", result.Error))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	follow := models.Follow{
		FollowerID: user.ID,
		Follower:   user,
		FolloweeID: followee.ID,
		Followee:   followee,
	}

	transactionErr := models.DB.Transaction(func(tx *gorm.DB) error {

		if createFollowResult := tx.Create(&follow); createFollowResult.Error != nil {
			fmt.Println(fmt.Errorf("[ERROR] %v", createFollowResult.Error))
			return createFollowResult.Error
		}

		followedNotifications := models.NewFollowedNotification(followee, follow)
		if createNotificationResult := tx.Create(&followedNotifications); createNotificationResult.Error != nil {
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

	c.JSON(http.StatusCreated, gin.H{"followed": follow})
}

func Unfollow(c *gin.Context) {
	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)

	followeeId, _ := strconv.Atoi(c.Param("id"))

	result := models.DB.Delete(&models.Follow{}, "follower_id = ? AND followee_id = ?", user.ID, followeeId)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"unfollowed": result})
}

func GetUsersFollowers(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	follows := []models.Follow{}
	models.DB.Preload("Follower").Where("followee_id = ?", userId).Find(&follows)

	followers := []models.User{}
	for _, follow := range follows {
		followers = append(followers, follow.Follower)
	}

	c.JSON(http.StatusOK, gin.H{"followers": followers})
}

func GetUsersFollowings(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	follows := []models.Follow{}
	models.DB.Preload("Followee").Where("follower_id = ?", userId).Find(&follows)

	followings := []models.User{}
	for _, follow := range follows {
		followings = append(followings, follow.Followee)
	}

	c.JSON(http.StatusOK, gin.H{"followings": followings})
}
