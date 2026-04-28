package restaurant

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mariorinawan12/food-api/internal/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(search string, categoryID uint, page, limit int) ([]domain.Restaurant, int64, error)
	FindByID(id uint) (*domain.Restaurant, error)
	FindByUserID(userID uint) ([]domain.Restaurant, error)
	Create(restaurant *domain.Restaurant) error
	Update(restaurant *domain.Restaurant) error
	Delete(id uint) error
	IsOwnedBy(restaurantID, userID uint) bool
}

type repository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewRepository(db *gorm.DB, rdb *redis.Client) Repository {
	return &repository{db, rdb}
}

var ctx = context.Background()

func (r *repository) FindAll(search string, categoryID uint, page, limit int) ([]domain.Restaurant, int64, error) {
	cacheKey := fmt.Sprintf("restaurants:all:%s:%d:%d:%d", search, categoryID, page, limit)

	cached, err := r.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var result struct {
			Data  []domain.Restaurant
			Total int64
		}
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result.Data, result.Total, nil
		}
	}

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
	if err := query.Offset(offset).Limit(limit).Find(&restaurants).Error; err != nil {
		return nil, 0, err
	}

	data := struct {
		Data  []domain.Restaurant
		Total int64
	}{restaurants, total}
	if b, err := json.Marshal(data); err == nil {
		r.rdb.Set(ctx, cacheKey, b, 5*time.Minute)
	}

	return restaurants, total, nil
}

func (r *repository) FindByID(id uint) (*domain.Restaurant, error) {
	cacheKey := fmt.Sprintf("restaurants:%d", id)

	cached, err := r.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var restaurant domain.Restaurant
		if err := json.Unmarshal([]byte(cached), &restaurant); err == nil {
			return &restaurant, nil
		}
	}

	var restaurant domain.Restaurant
	if err := r.db.Preload("Category").First(&restaurant, id).Error; err != nil {
		return nil, err
	}

	if b, err := json.Marshal(restaurant); err == nil {
		r.rdb.Set(ctx, cacheKey, b, 5*time.Minute)
	}

	return &restaurant, nil
}

func (r *repository) FindByUserID(userID uint) ([]domain.Restaurant, error) {
	var restaurants []domain.Restaurant
	err := r.db.Preload("Category").Where("created_by = ?", userID).Find(&restaurants).Error
	return restaurants, err
}

func (r *repository) Create(restaurant *domain.Restaurant) error {
	if err := r.db.Create(restaurant).Error; err != nil {
		return err
	}
	r.invalidateListCache()
	return r.db.Preload("Category").First(restaurant, restaurant.ID).Error
}

func (r *repository) Update(restaurant *domain.Restaurant) error {
	if err := r.db.Save(restaurant).Error; err != nil {
		return err
	}
	r.invalidateListCache()
	r.rdb.Del(ctx, fmt.Sprintf("restaurants:%d", restaurant.ID))
	return r.db.Preload("Category").First(restaurant, restaurant.ID).Error
}

func (r *repository) Delete(id uint) error {
	if err := r.db.Delete(&domain.Restaurant{}, id).Error; err != nil {
		return err
	}
	r.invalidateListCache()
	r.rdb.Del(ctx, fmt.Sprintf("restaurants:%d", id))
	return nil
}

func (r *repository) IsOwnedBy(restaurantID, userID uint) bool {
	var count int64
	r.db.Model(&domain.Restaurant{}).
		Where("id = ? AND created_by = ?", restaurantID, userID).
		Count(&count)
	return count > 0
}

func (r *repository) invalidateListCache() {
	keys, _ := r.rdb.Keys(ctx, "restaurants:all:*").Result()
	if len(keys) > 0 {
		r.rdb.Del(ctx, keys...)
	}
}
