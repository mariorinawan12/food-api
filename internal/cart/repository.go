package cart

import (
	"github.com/mariorinawan12/food-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	FindAllByUserID(userID uint) ([]domain.Cart, error)
	FindByID(id uint) (*domain.Cart, error)
	FindByUserIDAndRestaurantID(userID, restaurantID uint) (*domain.Cart, error)
	Create(cart *domain.Cart) error
	Delete(id uint) error
	FindCartItemByID(id uint) (*domain.CartItem, error)
	FindCartItemByMenuID(cartID, menuID uint) (*domain.CartItem, error)
	AddItem(item *domain.CartItem) error
	UpdateItem(item *domain.CartItem) error
	DeleteItem(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAllByUserID(userID uint) ([]domain.Cart, error) {
	var carts []domain.Cart
	err := r.db.Preload("Restaurant").Preload("CartItems.Menu").
		Where("user_id = ?", userID).Find(&carts).Error
	return carts, err
}

func (r *repository) FindByID(id uint) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.Preload("Restaurant").Preload("CartItems.Menu").
		First(&cart, id).Error
	return &cart, err
}

func (r *repository) FindByUserIDAndRestaurantID(userID uint, restaurantID uint) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.Preload("Restaurant").Preload("CartItems.Menu").
		Where("user_id = ? AND restaurant_id = ?", userID, restaurantID).
		First(&cart).Error
	return &cart, err
}

func (r *repository) Create(cart *domain.Cart) error {
	if err := r.db.Create(cart).Error; err != nil {
		return err
	}
	return r.db.Preload("Restaurant").Preload("CartItems.Menu").
		First(cart, cart.ID).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&domain.Cart{}, id).Error
}

func (r *repository) FindCartItemByID(id uint) (*domain.CartItem, error) {
	var item domain.CartItem
	err := r.db.Preload("Menu").First(&item, id).Error
	return &item, err
}

func (r *repository) FindCartItemByMenuID(cartID uint, menuID uint) (*domain.CartItem, error) {
	var item domain.CartItem
	err := r.db.Where("cart_id = ? AND menu_id = ?", cartID, menuID).
		Preload("Menu").First(&item).Error
	return &item, err
}

func (r *repository) AddItem(item *domain.CartItem) error {
	return r.db.Create(item).Error
}

func (r *repository) UpdateItem(item *domain.CartItem) error {
	return r.db.Save(item).Error
}

func (r *repository) DeleteItem(id uint) error {
	return r.db.Delete(&domain.CartItem{}, id).Error
}
