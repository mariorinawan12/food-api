package domain

import "time"

type Cart struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `json:"user_id"`
	User         User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	RestaurantID uint       `json:"restaurant_id"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant"`
	CartItems    []CartItem `gorm:"foreignKey:CartID" json:"cart_items"`
	CreatedAt    time.Time  `json:"created_at"`
}

type CartItem struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	CartID   uint `json:"cart_id"`
	MenuID   uint `json:"menu_id"`
	Menu     Menu `gorm:"foreignKey:MenuID" json:"menu"`
	Quantity int  `gorm:"not null" json:"quantity"`
}
