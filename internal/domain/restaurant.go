package domain

import "time"

type RestaurantCategory struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100;not null" json:"name"`
}

type Restaurant struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	Name        string             `gorm:"size:100;not null" json:"name"`
	Description string             `gorm:"type:text" json:"description"`
	Address     string             `gorm:"size:255" json:"address"`
	CategoryID  uint               `gorm:"category_id"`
	Category    RestaurantCategory `gorm:"foreignKey:CategoryID" json:"category"`
	CreatedBy   uint               `json:"created_by"`
	User        User               `gorm:"foreignKey:CreatedBy;constraint:OnDelete:Cascade" json:"created_by_user"`
	CreatedAt   time.Time          `json:"created_at"`
}
