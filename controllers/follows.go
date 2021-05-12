package controllers

import (
	"errors"
	"fmt"
	"gin-twitter/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FollowInput struct {
	FolloweeID uint `json:"followee_id" binding:"required"`
}

type UnfollowInput struct {
	FolloweeID uint `json:"followee_id" binding:"required"`
}

func Follow(c *gin.Context) {
	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)

	var input FollowInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println(fmt.Errorf("[ERROR]: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.ID == input.FolloweeID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Users may not follow themselves"})
		return
	}

	var followee models.User
	result := models.DB.First(&followee, input.FolloweeID)
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

	if result := models.DB.Create(&follow); result.Error != nil {
		fmt.Println(fmt.Errorf("[ERROR]: %v", result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"payload": follow})
}

func Unfollow(c *gin.Context) {
	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)

	var input UnfollowInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := models.DB.Delete(&models.Follow{}, "follower_id = ? AND followee_id = ?", user.ID, input.FolloweeID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"payload": result})
}
