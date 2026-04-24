package cart

import "time"

// request
type CreateCartRequest struct {
	RestaurantID uint `json:"restaurant_id" validate:"required"`
}

type AddItemRequest struct {
	MenuID   uint `json:"menu_id" validate:"required"`
	Quantity int  `json:"quantity" validate:"required,gt=0"`
}

type UpdateItemRequest struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
}

// response
type CartItemResponse struct {
	ID       uint    `json:"id"`
	MenuID   uint    `json:"menu_id"`
	MenuName string  `json:"menu_name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
}

type CartResponse struct {
	ID             uint               `json:"id"`
	UserID         uint               `json:"user_id"`
	RestaurantID   uint               `json:"restaurant_id"`
	RestaurantName string             `json:"restaurant_name"`
	Items          []CartItemResponse `json:"items"`
	TotalPrice     float64            `json:"total_price"`
	CreatedAt      time.Time          `json:"created_at"`
}
