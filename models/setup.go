package models

import (
	"log"
	"os"
	"strconv"
	"time"

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

	sqlDB, _ := database.DB()

	maxIdleConns, err := strconv.Atoi(os.Getenv("MAX_IDLE_CONNECTIONS"))
	if err != nil {
		log.Fatalf("MAX_IDLE_CONNECTIONS should be an integer, given %v", os.Getenv("MAX_IDLE_CONNECTIONS"))
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)

	maxOpenConns, err := strconv.Atoi(os.Getenv("MAX_OPEN_CONNECTIONS"))
	if err != nil {
		log.Fatalf("MAX_OPEN_CONNECTIONS should be an integer, given %v", os.Getenv("MAX_OPEN_CONNECTIONS"))
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)

	connMaxLifetime, err := strconv.Atoi(os.Getenv("CONNECTION_MAX_LIFETIME_MINUTES"))
	if err != nil {
		log.Fatalf("CONNECTION_MAX_LIFETIME_MINUTES should be an integer, given %v", os.Getenv("CONNECTION_MAX_LIFETIME_MINUTES"))
	}
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Minute)

	database.AutoMigrate(&User{})
	database.AutoMigrate(&Tweet{})
	database.AutoMigrate(&Follow{})
	database.AutoMigrate(&Like{})
	database.AutoMigrate(&Notification{})
	DB = database
}
