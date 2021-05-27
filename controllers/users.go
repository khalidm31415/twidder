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

// FindUsers godoc
// @Description get all users
// @Tags users
// @Success 200 {array} models.User
// @Router /users [get]
func FindUsers(c *gin.Context) {
	var users []models.User
	models.DB.Scopes(Paginate(c)).Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// FindUsers godoc
// @Description get user by id
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func FindUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	result := models.DB.Take(&user, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println(fmt.Errorf("[ERROR] %v", result.Error))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// FindUsers godoc
// @Description get a user's tweets
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {array} models.Tweet
// @Router /users/{id}/tweets [get]
func GetUsersTweets(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	var tweets []models.Tweet
	models.DB.Preload("User").Find(&tweets, "tweets.user_id = ?", userId)
	c.JSON(http.StatusOK, gin.H{"tweets": tweets})
}
