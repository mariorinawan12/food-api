package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mariorinawan12/food-api/internal/auth"
	"github.com/mariorinawan12/food-api/internal/cart"
	"github.com/mariorinawan12/food-api/internal/menu"
	"github.com/mariorinawan12/food-api/internal/middleware"
	"github.com/mariorinawan12/food-api/internal/order"
	"github.com/mariorinawan12/food-api/internal/restaurant"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	// Auth
	authRepo := auth.NewRepository(db)
	authUsecase := auth.NewUseCase(authRepo)
	authHandler := auth.NewHandler(authUsecase)
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

	// Restaurant
	restaurantRepo := restaurant.NewRepository(db)
	restaurantUsecase := restaurant.NewUsecase(restaurantRepo)
	restaurantHandler := restaurant.NewHandler(restaurantUsecase)

	// Menu
	menuRepo := menu.NewRepository(db)
	menuUsecase := menu.NewUsecase(menuRepo, restaurantRepo)
	menuHandler := menu.NewHandler(menuUsecase)

	// Cart
	cartRepo := cart.NewRepository(db)
	cartUsecase := cart.NewUsecase(cartRepo)
	cartHandler := cart.NewHandler(cartUsecase)

	// Order
	orderRepo := order.NewRepository(db)
	orderUsecase := order.NewUsecase(orderRepo, db)
	orderHandler := order.NewHandler(orderUsecase)

	// Public
	api.GET("/restaurants", restaurantHandler.GetAll)
	api.GET("/restaurants/:id", restaurantHandler.GetByID)
	api.GET("/menus", menuHandler.GetAll)
	api.GET("/menus/:id", menuHandler.GetByID)

	// User only
	user := api.Group("/")
	user.Use(middleware.AuthMiddleware(), middleware.UserOnly())
	user.POST("/cart", cartHandler.CreateCart)
	user.GET("/cart", cartHandler.GetAllCarts)
	user.GET("/cart/:id", cartHandler.GetCartByID)
	user.DELETE("/cart/:id", cartHandler.DeleteCart)
	user.POST("/cart/:id/items", cartHandler.AddItem)
	user.PUT("/cart/:id/items/:item_id", cartHandler.UpdateItem)
	user.DELETE("/cart/:id/items/:item_id", cartHandler.DeleteItem)
	user.POST("/checkout", orderHandler.Checkout)
	user.GET("/orders", orderHandler.GetMyOrders)
	user.GET("/orders/:id", orderHandler.GetByID)
	user.POST("/orders/:id/pay", orderHandler.Pay)
	user.POST("/orders/:id/cancel", orderHandler.Cancel)

	// Restaurant admin
	restoAdmin := api.Group("/")
	restoAdmin.Use(middleware.AuthMiddleware(), middleware.RestaurantAdminOnly())
	restoAdmin.GET("/my-restaurants", restaurantHandler.GetMyRestaurants)
	restoAdmin.POST("/restaurants", restaurantHandler.Create)
	restoAdmin.PUT("/restaurants/:id", restaurantHandler.Update)
	restoAdmin.DELETE("/restaurants/:id", restaurantHandler.Delete)
	restoAdmin.POST("/menus", menuHandler.Create)
	restoAdmin.PUT("/menus/:id", menuHandler.Update)
	restoAdmin.DELETE("/menus/:id", menuHandler.Delete)
	restoAdmin.GET("/resto-orders", orderHandler.GetRestaurantOrders)
	restoAdmin.POST("/orders/:id/process", orderHandler.Process)
	restoAdmin.POST("/orders/:id/deliver", orderHandler.Deliver)

	// Super admin
	superAdmin := api.Group("/admin")
	superAdmin.Use(middleware.AuthMiddleware(), middleware.SuperAdminOnly())
	superAdmin.GET("/orders", orderHandler.GetAllOrders)
	// superAdmin.GET("/users", authHandler.GetAllUsers)

	return r
}
