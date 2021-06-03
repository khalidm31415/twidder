package router

import (
	"twidder/controllers"
	"twidder/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	middelwares.InitAuthtMiddleware()

	r := gin.Default()

	r.GET("/ping", controllers.Ping)

	r.POST("/auth/signup", controllers.Signup)
	r.POST("/auth/login", controllers.Login)
	r.POST("/auth/logout", controllers.Logout)
	r.GET("/auth/current-user", middelwares.AuthMiddleware.MiddlewareFunc(), controllers.CurrentUser)
	r.DELETE("/auth/deactivate-account", middelwares.AuthMiddleware.MiddlewareFunc(), controllers.DeactivateAccount)
	r.POST("/auth/reactivate-account", controllers.ReactivateAccount)

	r.GET("/users", controllers.FindUsers)
	r.GET("/users/:id", controllers.FindUser)
	r.GET("/users/:id/tweets", controllers.GetUsersTweets)

	r.POST("/users/:id/follow", middelwares.AuthMiddleware.MiddlewareFunc(), controllers.Follow)
	r.POST("/users/:id/unfollow", middelwares.AuthMiddleware.MiddlewareFunc(), controllers.Unfollow)
	r.GET("/users/:id/followers", controllers.GetUsersFollowers)
	r.GET("/users/:id/followings", controllers.GetUsersFollowings)

	r.GET("/tweets", controllers.FindTweets)
	r.GET("/tweets/:id", controllers.FindTweet)
	r.POST("/tweets", middelwares.AuthMiddleware.MiddlewareFunc(), controllers.CreateTweet)
	r.DELETE("/tweets/:id", middelwares.AuthMiddleware.MiddlewareFunc(), controllers.DeleteTweet)
	r.GET("/tweets/timeline", middelwares.AuthMiddleware.MiddlewareFunc(), controllers.Timeline)

	r.POST("/tweets/:id/like", middelwares.AuthMiddleware.MiddlewareFunc(), controllers.Like)
	r.POST("/tweets/:id/unlike", middelwares.AuthMiddleware.MiddlewareFunc(), controllers.Unlike)
	r.GET("/tweets/:id/likes", controllers.Likes)

	r.POST("/tweets/:id/reply", middelwares.AuthMiddleware.MiddlewareFunc(), controllers.Reply)
	r.GET("/tweets/:id/replies", controllers.Replies)

	r.GET("/notifications", middelwares.AuthMiddleware.MiddlewareFunc(), controllers.Notifications)

	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return r
}
