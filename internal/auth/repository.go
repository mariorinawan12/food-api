package auth

import (
	"github.com/mariorinawan12/food-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(email string) (*domain.User, error)
	Create(user *domain.User) error
	FindRoleByName(name string) (*domain.Role, error)
	FindAll(page int, limit int) ([]domain.User, int64, error)
	FindByID(id uint) (*domain.User, error)
	UpdatePassword(userID uint, hashedPassword string) error
	UpdateProfile(userID uint, name string, email string) error
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

func (r *repository) FindAll(page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	query := r.db.Model(&domain.User{}).Preload("Role")
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *repository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Role").First(&user, id).Error
	return &user, err
}

func (r *repository) UpdatePassword(userID uint, hashedPassword string) error {
	return r.db.Model(&domain.User{}).
		Where("id =?", userID).
		Update("password", hashedPassword).Error
}

func (r *repository) UpdateProfile(userID uint, name string, email string) error {
	return r.db.Model(&domain.User{}).
		Where("id = ?", userID).
		Updates(&domain.User{
			Name:  name,
			Email: email,
		}).Error
}
