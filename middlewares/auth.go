package middelwares

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"
	"twidder/models"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var identityKey = "id"
var AuthMiddleware *jwt.GinJWTMiddleware

type LoginInput struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func InitAuthtMiddleware() {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(os.Getenv("JWT_SECRET_KEY")),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals LoginInput
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			givenUsername := loginVals.Username
			givenPassword := loginVals.Password

			var user models.User
			result := models.DB.Where(&models.User{Username: givenUsername}).Take(&user)

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, jwt.ErrFailedAuthentication
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(givenPassword)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return user.ID, nil

		},
		PayloadFunc: func(userID interface{}) jwt.MapClaims {
			if userID != nil {
				return jwt.MapClaims{
					identityKey: userID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			userID := claims[identityKey]
			var user models.User
			models.DB.First(&user, userID)
			return user
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(models.User); ok {
				return true
			}
			return false
		},
		TimeFunc:       time.Now,
		SendCookie:     true,
		SecureCookie:   false,
		CookieHTTPOnly: true,
		CookieDomain:   "localhost",
		CookieName:     "token",
		TokenLookup:    "cookie:token",
		CookieSameSite: http.SameSiteDefaultMode,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	AuthMiddleware = authMiddleware
}
