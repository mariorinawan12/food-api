package restaurant

import (
	"github.com/mariorinawan12/food-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(search string, categoryID uint, page, limit int) ([]domain.Restaurant, int64, error)
	FindByID(id uint) (*domain.Restaurant, error)
	Create(restaurant *domain.Restaurant) error
	Update(restaurant *domain.Restaurant) error
	Delete(id uint) error
	IsOwnedBy(restaurantID, userID uint) bool
	FindByUserID(userID uint) ([]domain.Restaurant, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAll(search string, categoryID uint, page int, limit int) ([]domain.Restaurant, int64, error) {
	var restaurants []domain.Restaurant
	var total int64

	query := r.db.Model(&domain.Restaurant{}).Preload("Category")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	if categoryID != 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&restaurants).Error
	return restaurants, total, err
}

func (r *repository) FindByID(id uint) (*domain.Restaurant, error) {
	var restaurant domain.Restaurant
	err := r.db.Preload("Category").First(&restaurant, id).Error
	return &restaurant, err
}

func (r *repository) Create(restaurant *domain.Restaurant) error {
	if err := r.db.Create(restaurant).Error; err != nil {
		return err
	}
	return r.db.Preload("Category").First(restaurant, restaurant.ID).Error
}

func (r *repository) Update(restaurant *domain.Restaurant) error {
	if err := r.db.Save(restaurant).Error; err != nil {
		return err
	}
	return r.db.Preload("Category").First(restaurant, restaurant.ID).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&domain.Restaurant{}, id).Error
}

func (r *repository) IsOwnedBy(restaurantID uint, userID uint) bool {
	var count int64
	r.db.Model(&domain.Restaurant{}).
		Where("id = ? AND created_by = ?", restaurantID, userID).
		Count(&count)
	return count > 0
}

func (r *repository) FindByUserID(userID uint) ([]domain.Restaurant, error) {
	var restaurants []domain.Restaurant
	err := r.db.Preload("Category").Where("created_by = ?", userID).Find(&restaurants).Error

	return restaurants, err
}
