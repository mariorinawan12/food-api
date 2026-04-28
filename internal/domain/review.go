package domain

import "time"

type Review struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `json:"user_id"`
	User         User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	RestaurantID uint       `json:"restaurant_id"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID;constraint:OnDelete:CASCADE" json:"restaurant"`
	Rating       int        `gorm:"not null" json:"rating"`
	Comment      string     `gorm:"type:text" json:"comment"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
