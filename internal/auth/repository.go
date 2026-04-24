package auth

import (
	"github.com/mariorinawan12/food-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(email string) (*domain.User, error)
	Create(user *domain.User) error
	FindRoleByName(name string) (*domain.Role, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *repository) Create(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}

	return r.db.Preload("Role").First(user, user.ID).Error
}

func (r *repository) FindRoleByName(name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Where("name =?", name).First(&role).Error
	return &role, err
}
