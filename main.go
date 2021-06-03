package main

import (
	"fmt"
	"os"
	"twidder/middlewares"
	"twidder/models"
	"twidder/router"

	"github.com/joho/godotenv"

	_ "twidder/docs"
)

// @title Twidder
// @description Imitating twitter backend API with Gin, GORM, and MySQL.

// @host localhost:8080
func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
	}

	models.ConnectDatabase()
	middelwares.InitAuthtMiddleware()

	r := router.SetupRouter()

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
