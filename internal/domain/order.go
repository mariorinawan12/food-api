package domain

import "time"

type Order struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	UserID     uint        `json:"user_id"`
	User       User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Address    string      `gorm:"size:255;not null" json:"address"`
	Status     string      `gorm:"size:50;default:payment_pending" json:"status"`
	TotalPrice float64     `gorm:"not null" json:"total_price"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
	CreatedAt  time.Time   `json:"created_at"`
}

type OrderItem struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	OrderID  uint    `json:"order_id"`
	MenuID   uint    `json:"menu_id"`
	Menu     Menu    `gorm:"foreignKey:MenuID" json:"menu"`
	Quantity int     `gorm:"not null" json:"quantity"`
	Price    float64 `gorm:"not null" json:"price"`
}
