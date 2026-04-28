package order

import "time"

type CheckoutRequest struct {
	CartID  uint   `json:"-"`
	Address string `json:"address" validate:"required"`
}

type OrderItemResponse struct {
	ID       uint    `json:"id"`
	MenuID   uint    `json:"menu_id"`
	MenuName string  `json:"menu_name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Subtotal float64 `json:"subtotal"`
}

type OrderResponse struct {
	ID         uint                `json:"id"`
	UserID     uint                `json:"user_id"`
	Address    string              `json:"address"`
	Status     string              `json:"status"`
	TotalPrice float64             `json:"total_price"`
	Items      []OrderItemResponse `json:"items"`
	CreatedAt  time.Time           `json:"created_at"`
}
