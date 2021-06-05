package test

import (
	"testing"
	"twidder/models"
	"twidder/router"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	mock   sqlmock.Sqlmock
	router *gin.Engine
}

func (suite *Suite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	router := router.SetupRouter()
	suite.router = router
}

func (suite *Suite) SetupTest() {
	sqlDB, mock, err := sqlmock.New()
	suite.NoError(err)
	suite.mock = mock

	database, _ := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	models.DB = database
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
