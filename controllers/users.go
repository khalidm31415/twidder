package controllers

import (
	"fmt"
	"gin-twitter/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var identityKey = "id"

type CreateUserInput struct {
	Username    string `binding:"required"`
	DisplayName string `binding:"required"`
	Email       string `binding:"required"`
	Password    string `binding:"required"`
}

type ReactivateAccountInput struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

func FindUsers(c *gin.Context) {
	var users []models.User
	models.DB.Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

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

func GetUsersTweets(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	var tweets []models.Tweet
	models.DB.Preload("User").Find(&tweets, "tweets.user_id = ?", userId)
	c.JSON(http.StatusOK, gin.H{"tweets": tweets})
}

func Signup(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	user := models.User{
		Username:    input.Username,
		DisplayName: input.DisplayName,
		Email:       input.Email,
		Password:    string(hashedPassword),
	}

	if dbc := models.DB.Create(&user); dbc.Error != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", dbc.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func CurrentUser(c *gin.Context) {
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   user.(models.User).ID,
		"username": user.(models.User).Username,
	})
}

func DeactivateAccount(c *gin.Context) {
	v, _ := c.Get(identityKey)
	user, _ := v.(models.User)

	transactionErr := models.DB.Transaction(func(tx *gorm.DB) error {
		deleteLikesResult := models.DB.Delete(&models.Like{}, "likes.user_id = ?", user.ID)
		if deleteLikesResult.Error != nil {
			fmt.Println(fmt.Errorf("[ERROR] %v", deleteLikesResult.Error))
			return deleteLikesResult.Error
		}

		deleteFollowsResult := models.DB.Delete(&models.Follow{}, "follows.follower_id = ? OR follows.followee_id", user.ID, user.ID)
		if deleteFollowsResult.Error != nil {
			fmt.Println(fmt.Errorf("[ERROR] %v", deleteFollowsResult.Error))
			return deleteFollowsResult.Error
		}

		deleteTweetsResult := models.DB.Delete(&models.Tweet{}, "tweets.user_id = ?", user.ID)
		if deleteTweetsResult.Error != nil {
			fmt.Println(fmt.Errorf("[ERROR] %v", deleteTweetsResult.Error))
			return deleteTweetsResult.Error
		}

		deleteUserResult := models.DB.Delete(&user)
		if deleteUserResult.Error != nil {
			fmt.Println(fmt.Errorf("[ERROR] %v", deleteUserResult.Error))
			return deleteUserResult.Error
		}
		return nil
	})

	if transactionErr != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", transactionErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"deleted": user})
}

func ReactivateAccount(c *gin.Context) {
	var input ReactivateAccountInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	givenUsername := input.Username
	givenPassword := input.Password

	var user models.User
	result := models.DB.Unscoped().Where(&models.User{Username: givenUsername}).Take(&user)
	if result.Error != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", result.Error))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(givenPassword)); err != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", result.Error))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect credentials"})
		return
	}

	transactionErr := models.DB.Transaction(func(tx *gorm.DB) error {
		restoreLikesResult := models.DB.Model(&models.Like{}).Unscoped().Where("likes.user_id = ?", user.ID)
		if restoreLikesResult.Error != nil {
			fmt.Println(fmt.Errorf("[ERROR] %v", restoreLikesResult.Error))
			return restoreLikesResult.Error
		}

		restoreFollowsResult := models.DB.Model(&models.Follow{}).Unscoped().Where("follows.follower_id = ? OR follows.followee_id", user.ID, user.ID)
		if restoreFollowsResult.Error != nil {
			fmt.Println(fmt.Errorf("[ERROR] %v", restoreFollowsResult.Error))
			return restoreFollowsResult.Error
		}

		restoreTweetsResult := models.DB.Model(&models.Tweet{}).Unscoped().Where("tweets.user_id = ?", user.ID)
		if restoreTweetsResult.Error != nil {
			fmt.Println(fmt.Errorf("[ERROR] %v", restoreTweetsResult.Error))
			return restoreTweetsResult.Error
		}

		restoreUserResult := models.DB.Model(&user).Update("deleted_at", nil)
		if restoreUserResult.Error != nil {
			fmt.Println(fmt.Errorf("[ERROR] %v", result.Error))
			return restoreUserResult.Error
		}

		return nil
	})

	if transactionErr != nil {
		fmt.Println(fmt.Errorf("[ERROR] %v", transactionErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
