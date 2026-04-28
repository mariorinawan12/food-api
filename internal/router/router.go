package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/mariorinawan12/food-api/internal/auth"
	"github.com/mariorinawan12/food-api/internal/cart"
	"github.com/mariorinawan12/food-api/internal/menu"
	"github.com/mariorinawan12/food-api/internal/middleware"
	"github.com/mariorinawan12/food-api/internal/order"
	"github.com/mariorinawan12/food-api/internal/restaurant"
	"github.com/mariorinawan12/food-api/internal/review"
)

func SetupRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	// Auth
	authRepo := auth.NewRepository(db)
	authUsecase := auth.NewUseCase(authRepo)
	authHandler := auth.NewHandler(authUsecase)
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

	authUser := api.Group("/")
	authUser.Use(middleware.AuthMiddleware())
	authUser.PUT("/change-password", authHandler.ChangePassword)
	authUser.PUT("/update-profile", authHandler.UpdateProfile)

	// Restaurant
	restaurantRepo := restaurant.NewRepository(db, rdb)
	restaurantUsecase := restaurant.NewUsecase(restaurantRepo)
	restaurantHandler := restaurant.NewHandler(restaurantUsecase)

	// Menu
	menuRepo := menu.NewRepository(db, rdb)
	menuUsecase := menu.NewUsecase(menuRepo, restaurantRepo)
	menuHandler := menu.NewHandler(menuUsecase)

	// Cart
	cartRepo := cart.NewRepository(db)
	cartUsecase := cart.NewUsecase(cartRepo, menuRepo)
	cartHandler := cart.NewHandler(cartUsecase)

	// Order
	orderRepo := order.NewRepository(db)
	orderUsecase := order.NewUsecase(orderRepo, cartRepo, db)
	orderHandler := order.NewHandler(orderUsecase)

	// Review
	reviewRepo := review.NewRepository(db, rdb)
	reviewUsecase := review.NewUsecase(reviewRepo)
	reviewHandler := review.NewHandler(reviewUsecase)

	// Public
	api.GET("/restaurants", restaurantHandler.GetAll)
	api.GET("/restaurants/:restaurant_id", restaurantHandler.GetByID)
	api.GET("/restaurants/:restaurant_id/menus", menuHandler.GetByRestaurant)
	api.GET("/restaurants/:restaurant_id/reviews", reviewHandler.GetAllByRestaurant)
	api.GET("/menus", menuHandler.GetAll)
	api.GET("/menus/:menu_id", menuHandler.GetByID)

	// User only
	user := api.Group("/")
	user.Use(middleware.AuthMiddleware(), middleware.UserOnly())
	user.POST("/restaurants/:restaurant_id/cart", cartHandler.CreateCart)
	user.GET("/cart", cartHandler.GetAllCarts)
	user.GET("/cart/:cart_id", cartHandler.GetCartByID)
	user.DELETE("/cart/:cart_id", cartHandler.DeleteCart)
	user.POST("/cart/:cart_id/items", cartHandler.AddItem)
	user.PUT("/cart/:cart_id/items/:item_id", cartHandler.UpdateItem)
	user.DELETE("/cart/:cart_id/items/:item_id", cartHandler.DeleteItem)
	user.POST("/cart/:cart_id/checkout", orderHandler.Checkout)
	user.GET("/orders", orderHandler.GetMyOrders)
	user.GET("/orders/:order_id", orderHandler.GetByID)
	user.POST("/orders/:order_id/pay", orderHandler.Pay)
	user.POST("/orders/:order_id/cancel", orderHandler.Cancel)
	user.POST("/restaurants/:restaurant_id/reviews", reviewHandler.Create)
	user.PUT("/reviews/:review_id", reviewHandler.Update)
	user.DELETE("/reviews/:review_id", reviewHandler.Delete)

	// Restaurant admin
	restoAdmin := api.Group("/")
	restoAdmin.Use(middleware.AuthMiddleware(), middleware.RestaurantAdminOnly())
	restoAdmin.GET("/my-restaurants", restaurantHandler.GetMyRestaurants)
	restoAdmin.POST("/restaurants", restaurantHandler.Create)
	restoAdmin.PUT("/restaurants/:restaurant_id", restaurantHandler.Update)
	restoAdmin.DELETE("/restaurants/:restaurant_id", restaurantHandler.Delete)
	restoAdmin.POST("/restaurants/:restaurant_id/menus", menuHandler.Create)
	restoAdmin.PUT("/menus/:menu_id", menuHandler.Update)
	restoAdmin.DELETE("/menus/:menu_id", menuHandler.Delete)
	restoAdmin.GET("/resto-orders", orderHandler.GetRestaurantOrders)
	restoAdmin.GET("/resto-orders/:order_id", orderHandler.GetOrderDetailByAdmin)
	restoAdmin.POST("/orders/:order_id/process", orderHandler.Process)
	restoAdmin.POST("/orders/:order_id/deliver", orderHandler.Deliver)

	// Super admin
	superAdmin := api.Group("/admin")
	superAdmin.Use(middleware.AuthMiddleware(), middleware.SuperAdminOnly())
	superAdmin.GET("/users", authHandler.GetAllUsers)
	superAdmin.GET("/orders", orderHandler.GetAllOrders)
	superAdmin.GET("/orders/:order_id", orderHandler.GetOrderDetail)
	return r
}
