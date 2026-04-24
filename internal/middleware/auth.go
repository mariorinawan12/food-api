package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mariorinawan12/food-api/internal/helper"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer") {
			helper.Error(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := helper.ValidateToken(tokenStr)
		if err != nil {
			helper.Error(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "super_admin" {
			helper.Error(c, http.StatusForbidden, "super admin only")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RestaurantAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "restaurant_admin" && role != "super_admin" {
			helper.Error(c, http.StatusForbidden, "restaurant admin only")
			c.Abort()
			return
		}
		c.Next()
	}
}

func UserOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "user" {
			helper.Error(c, http.StatusForbidden, "user only")
			c.Abort()
			return
		}
		c.Next()
	}
}
