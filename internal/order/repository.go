package order

import (
	"github.com/mariorinawan12/food-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	FindAllByUserID(userID uint, status string, page int, limit int) ([]domain.Order, int64, error)
	FindAll(status string, page int, limit int) ([]domain.Order, int64, error)
	FindAllByRestaurantAdminID(adminID uint, status string, page int, limit int) ([]domain.Order, int64, error)
	FindByID(id uint) (*domain.Order, error)
	Create(order *domain.Order) error
	UpdateStatus(id uint, status string) error
	IsOrderOwnedByAdmin(adminID uint, orderID uint) bool
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAllByUserID(userID uint, status string, page, limit int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	query := r.db.Model(&domain.Order{}).
		Preload("OrderItems.Menu").
		Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&orders).Error
	return orders, total, err
}

func (r *repository) FindAll(status string, page int, limit int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	query := r.db.Model(&domain.Order{}).
		Preload("OrderItems.Menu").
		Preload("User")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&orders).Error
	return orders, total, err
}

func (r *repository) FindAllByRestaurantAdminID(adminID uint, status string, page int, limit int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	query := r.db.Model(&domain.Order{}).
		Preload("OrderItems.Menu").
		Preload("User").
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Joins("JOIN menus ON menus.id = order_items.menu_id").
		Joins("JOIN restaurants ON restaurants.id = menus.restaurant_id").
		Where("restaurants.created_by = ?", adminID).
		Distinct("orders.id")

	if status != "" {
		query = query.Where("orders.status = ?", status)
	}

	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("orders.created_at desc").Find(&orders).Error
	return orders, total, err
}

func (r *repository) FindByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("OrderItems.Menu").Preload("User").First(&order.ID).Error
	return &order, err
}

func (r *repository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *repository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&domain.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *repository) IsOrderOwnedByAdmin(adminID uint, orderID uint) bool {
	var count int64
	r.db.Model(&domain.Restaurant{}).
		Joins("JOIN menus ON menus.restaurant_id = restaurants.id").
		Joins("JOIN order_items ON order_items.menu_id = menus.id").
		Where("order_items.order_id = ? AND restaurants.created_by = ?", orderID, adminID).
		Count(&count)
	return count > 0
}
