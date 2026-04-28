package config

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/mariorinawan12/food-api/internal/domain"
	"github.com/mariorinawan12/food-api/internal/helper"
)

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	return db
}

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})
	return rdb
}

func RunMigration(db *gorm.DB) {
	db.AutoMigrate(
		&domain.Role{},
		&domain.User{},
		&domain.RestaurantCategory{},
		&domain.Restaurant{},
		&domain.Menu{},
		&domain.Cart{},
		&domain.CartItem{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Review{},
	)

	// seed role default
	roles := []string{"super_admin", "restaurant_admin", "user"}
	for _, name := range roles {
		db.FirstOrCreate(&domain.Role{}, domain.Role{Name: name})
	}

	// seed super_admin account
	password, _ := helper.HashPassword("password")
	db.FirstOrCreate(&domain.User{}, domain.User{
		Name:     "Super Admin",
		Email:    "superadmin@mail.com",
		Password: password,
		RoleID:   1,
	})

	categories := []string{"Western Food", "Asian Food", "Local Food", "Snack", "Drink"}
	for _, name := range categories {
		db.FirstOrCreate(&domain.RestaurantCategory{}, domain.RestaurantCategory{Name: name})
	}

}
