package domain

import "time"

type Role struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:50;not null" json:"name"`
}

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	RoleID    uint      `json:"-"`
	Role      Role      `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
