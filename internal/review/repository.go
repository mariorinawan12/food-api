package review

import (
	"context"
	"fmt"

	"github.com/mariorinawan12/food-api/internal/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	FindAllByRestaurantID(restaurantID uint, page, limit int) ([]domain.Review, int64, error)
	FindByID(id uint) (*domain.Review, error)
	FindByUserAndRestaurant(userID, restaurantID uint) (*domain.Review, error)
	HasDeliveredOrder(userID, restaurantID uint) bool
	Create(review *domain.Review) error
	Update(review *domain.Review) error
	Delete(id uint) error
	UpdateRestaurantRating(restaurantID uint) error
}

type repository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewRepository(db *gorm.DB, rdb *redis.Client) Repository {
	return &repository{db, rdb}
}

var ctx = context.Background()

func (r *repository) FindAllByRestaurantID(restaurantID uint, page, limit int) ([]domain.Review, int64, error) {
	var reviews []domain.Review
	var total int64

	query := r.db.Model(&domain.Review{}).
		Preload("User").
		Where("restaurant_id = ?", restaurantID)

	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&reviews).Error
	return reviews, total, err
}

func (r *repository) FindByID(id uint) (*domain.Review, error) {
	var review domain.Review
	err := r.db.Preload("User").Preload("Restaurant").First(&review, id).Error
	return &review, err
}

func (r *repository) FindByUserAndRestaurant(userID, restaurantID uint) (*domain.Review, error) {
	var review domain.Review
	err := r.db.Where("user_id = ? AND restaurant_id = ?", userID, restaurantID).First(&review).Error
	return &review, err
}

func (r *repository) HasDeliveredOrder(userID, restaurantID uint) bool {
	var count int64
	r.db.Model(&domain.Order{}).
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Joins("JOIN menus ON menus.id = order_items.menu_id").
		Where("orders.user_id = ? AND menus.restaurant_id = ? AND orders.status = ?", userID, restaurantID, "delivered").
		Count(&count)
	return count > 0
}

func (r *repository) Create(review *domain.Review) error {
	if err := r.db.Create(review).Error; err != nil {
		return err
	}
	if err := r.UpdateRestaurantRating(review.RestaurantID); err != nil {
		return err
	}
	r.invalidateRestaurantCache(review.RestaurantID)
	return nil
}

func (r *repository) Update(review *domain.Review) error {
	if err := r.db.Save(review).Error; err != nil {
		return err
	}
	if err := r.UpdateRestaurantRating(review.RestaurantID); err != nil {
		return err
	}
	r.invalidateRestaurantCache(review.RestaurantID)
	return nil
}

func (r *repository) Delete(id uint) error {
	review, err := r.FindByID(id)
	if err != nil {
		return err
	}
	restaurantID := review.RestaurantID
	if err := r.db.Delete(&domain.Review{}, id).Error; err != nil {
		return err
	}
	if err := r.UpdateRestaurantRating(restaurantID); err != nil {
		return err
	}
	r.invalidateRestaurantCache(restaurantID)
	return nil
}

func (r *repository) UpdateRestaurantRating(restaurantID uint) error {
	var avg float64
	r.db.Model(&domain.Review{}).
		Where("restaurant_id = ?", restaurantID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avg)

	return r.db.Model(&domain.Restaurant{}).
		Where("id = ?", restaurantID).
		Update("average_rating", avg).Error
}

func (r *repository) invalidateRestaurantCache(restaurantID uint) {
	r.rdb.Del(ctx, fmt.Sprintf("restaurants:%d", restaurantID))
	keys, _ := r.rdb.Keys(ctx, "restaurants:all:*").Result()
	if len(keys) > 0 {
		r.rdb.Del(ctx, keys...)
	}
}
