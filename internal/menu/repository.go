package menu

import (
	"github.com/mariorinawan12/food-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(restaurantID uint, search string, page, limit int) ([]domain.Menu, int64, error)
	FindByID(id uint) (*domain.Menu, error)
	Create(menu *domain.Menu) error
	Update(menu *domain.Menu) error
	Delete(id uint) error
	IsOwnedBy(menuID, userID uint) bool
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAll(restaurantID uint, search string, page, limit int) ([]domain.Menu, int64, error) {
	var menus []domain.Menu
	var total int64

	query := r.db.Model(&domain.Menu{}).Preload("Restaurant")

	if restaurantID != 0 {
		query = query.Where("restaurant_id = ?", restaurantID)
	}
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&menus).Error
	return menus, total, err
}

func (r *repository) FindByID(id uint) (*domain.Menu, error) {
	var menu domain.Menu
	err := r.db.Preload("Restaurant").First(&menu, id).Error
	return &menu, err
}

func (r *repository) Create(menu *domain.Menu) error {
	return r.db.Create(menu).Error
}

func (r *repository) Update(menu *domain.Menu) error {
	if err := r.db.Save(menu).Error; err != nil {
		return err
	}
	return r.db.Preload("Restaurant").First(menu, menu.ID).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&domain.Menu{}, id).Error
}

func (r *repository) IsOwnedBy(menuID, userID uint) bool {
	var count int64
	r.db.Model(&domain.Menu{}).
		Joins("JOIN restaurants ON restaurants.id = menus.restaurant_id").
		Where("menus.id = ? AND restaurants.created_by = ?", menuID, userID).
		Count(&count)
	return count > 0
}
