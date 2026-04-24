package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mariorinawan12/food-api/config"
	"github.com/mariorinawan12/food-api/internal/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("JWT_SECRET:", os.Getenv("JWT_SECRET"))

	db := config.InitDB()
	config.RunMigration(db)

	r := router.SetupRouter(db)
	r.Run(":" + os.Getenv("APP_PORT"))
}
