package restaurant

import "time"

type CreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Address     string `json:"address" validate:"required"`
	CategoryID  uint   `json:"category_id" validate:"required"`
}

type UpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	CategoryID  uint   `json:"category_id"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type RestaurantResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Address      string    `json:"address"`
	CategoryName string    `json:"category_name"`
	CreatedBy    uint      `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}
