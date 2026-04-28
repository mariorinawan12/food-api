package menu

import "time"

type CreateRequest struct {
	RestaurantID uint    `json:"-"`
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description"`
	Price        float64 `json:"price" validate:"required,gt=0"`
	Category     string  `json:"category" validate:"required"`
}

type UpdateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

type MenuResponse struct {
	ID             uint      `json:"id"`
	RestaurantID   uint      `json:"restaurant_id"`
	RestaurantName string    `json:"restaurant_name"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          float64   `json:"price"`
	Category       string    `json:"category"`
	CreatedAt      time.Time `json:"created_at"`
}
