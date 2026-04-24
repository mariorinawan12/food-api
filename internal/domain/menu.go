package domain

import "time"

type Menu struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	RestaurantID uint       `json:"restaurant_id"`
	Restaurant   Restaurant `gorm:"primaryKey:RestaurantID" json:"restaurant"`
	Name         string     `gorm:"size:100;not null" json:"name"`
	Description  string     `gorm:"type:text" json:"description"`
	Price        float64    `gorm:"not null" json:"price"`
	Category     string     `gorm:"size:100" json:"category"`
	CreatedAt    time.Time  `json:"created_at"`
}
