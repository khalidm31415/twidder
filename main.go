package main

import (
	"fmt"
	"gin-twitter/controllers"
	"gin-twitter/middlewares"
	"gin-twitter/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gin-twitter/docs"
)

// @title Twitter with Gin and GORM
// @description A twitter-like API implemented with Gin, GORM, and MySQL.

// @host localhost:8080
func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	models.ConnectDatabase()

	r := gin.Default()

	authMiddleware := middelwares.GeAuthtMiddleware()

	r.GET("/ping", ping)

	r.POST("/users/signup", controllers.Signup)
	r.POST("/users//login", authMiddleware.LoginHandler)
	r.POST("/users//logout", authMiddleware.LogoutHandler)
	r.GET("/users/current-user", authMiddleware.MiddlewareFunc(), controllers.CurrentUser)
	r.DELETE("/users/deactivate-account", authMiddleware.MiddlewareFunc(), controllers.DeactivateAccount)
	r.POST("/users/reactivate-account", controllers.ReactivateAccount)

	r.GET("/users", controllers.FindUsers)
	r.GET("/users/:id", controllers.FindUser)
	r.GET("/users/:id/tweets", controllers.GetUsersTweets)
	r.GET("/users/:id/followers", controllers.GetUsersFollowers)
	r.GET("/users/:id/followings", controllers.GetUsersFollowings)

	r.GET("/tweets", controllers.FindTweets)
	r.GET("/tweets/:id", controllers.FindTweet)
	r.POST("/tweets", authMiddleware.MiddlewareFunc(), controllers.CreateTweet)
	r.DELETE("/tweets/:id", authMiddleware.MiddlewareFunc(), controllers.DeleteTweet)

	r.POST("/follow", authMiddleware.MiddlewareFunc(), controllers.Follow)
	r.POST("/unfollow", authMiddleware.MiddlewareFunc(), controllers.Unfollow)

	r.POST("/like", authMiddleware.MiddlewareFunc(), controllers.Like)
	r.POST("/unlike", authMiddleware.MiddlewareFunc(), controllers.Unlike)

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

// Ping godoc
// @Description ping the server.
// @Tags root
// @Success 200 {string} string "pong"
// @Router /ping [get]
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
