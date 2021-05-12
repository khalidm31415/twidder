package models

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("MYSQL_DSN")
	database, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	database.AutoMigrate(&User{})
	database.AutoMigrate(&Tweet{})
	database.AutoMigrate(&Follow{})
	database.AutoMigrate(&Like{})
	DB = database
}
