package controllers

import (
	"fmt"
	"net/http"
	"twidder/middlewares"
	"twidder/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var identityKey = "id"

type SignupInput struct {
	Username    string `binding:"required"`
	DisplayName string `binding:"required"`
	Email       string `binding:"required"`
	Password    string `binding:"required"`
}

type ReactivateAccountInput struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

// Signup godoc
// @Description create new user
// @Tags auth
// @Accept json
// @Produce json
// @Param signup body SignupInput true "Signup to create a new user"
// @Success 200 {object} models.User
// @Router /auth/signup [post]
func Signup(c *gin.Context) {
	var input SignupInput
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

// Login godoc
// @Description login
// @Tags auth
// @Accept json
// @Produce json
// @Param login body middelwares.LoginInput true "Login"
// @Success 200
// @Router /auth/login [post]
func Login(c *gin.Context) {
	middelwares.AuthMiddleware.LoginHandler(c)
}

// Logout godoc
// @Description logout
// @Tags auth
// @Accept json
// @Produce json
// @Success 200
// @Router /auth/logout [post]
func Logout(c *gin.Context) {
	middelwares.AuthMiddleware.LogoutHandler(c)
}

// CurrentUser godoc
// @Description check currently logged in user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200
// @Router /auth/current-user [get]
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
