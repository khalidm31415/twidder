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
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	r := gin.Default()

	authMiddleware := middelwares.GeAuthtMiddleware()

	r.GET("/ping", func(c *gin.Context) {
		fmt.Println(c.Request.Header)
		c.JSON(http.StatusOK, "pong")
	})

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
	r.POST("/tweets", authMiddleware.MiddlewareFunc(), controllers.CreateTweet)

	r.POST("/follow", authMiddleware.MiddlewareFunc(), controllers.Follow)
	r.POST("/unfollow", authMiddleware.MiddlewareFunc(), controllers.Unfollow)

	r.POST("/like", authMiddleware.MiddlewareFunc(), controllers.Like)
	r.POST("/unlike", authMiddleware.MiddlewareFunc(), controllers.Unlike)

	models.ConnectDatabase()

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
