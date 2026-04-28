package review

import "time"

type CreateRequest struct {
	RestaurantID uint   `json:"-"`
	Rating       int    `json:"rating" validate:"required,min=1,max=5"`
	Comment      string `json:"comment" validate:"required"`
}

type UpdateRequest struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment" validate:"required"`
}

type ReviewResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	UserName     string    `json:"user_name"`
	RestaurantID uint      `json:"restaurant_id"`
	Rating       int       `json:"rating"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
