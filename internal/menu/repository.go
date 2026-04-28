package menu

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
	FindAll(restaurantID uint, search string, page, limit int) ([]domain.Menu, int64, error)
	FindByID(id uint) (*domain.Menu, error)
	Create(menu *domain.Menu) error
	Update(menu *domain.Menu) error
	Delete(id uint) error
	IsOwnedBy(menuID, userID uint) bool
}

type repository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewRepository(db *gorm.DB, rdb *redis.Client) Repository {
	return &repository{db, rdb}
}

var ctx = context.Background()

func (r *repository) FindAll(restaurantID uint, search string, page, limit int) ([]domain.Menu, int64, error) {
	cacheKey := fmt.Sprintf("menus:all:%d:%s:%d:%d", restaurantID, search, page, limit)

	cached, err := r.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var result struct {
			Data  []domain.Menu
			Total int64
		}
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result.Data, result.Total, nil
		}
	}

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
	if err := query.Offset(offset).Limit(limit).Find(&menus).Error; err != nil {
		return nil, 0, err
	}

	data := struct {
		Data  []domain.Menu
		Total int64
	}{menus, total}
	if b, err := json.Marshal(data); err == nil {
		r.rdb.Set(ctx, cacheKey, b, 5*time.Minute)
	}

	return menus, total, nil
}

func (r *repository) FindByID(id uint) (*domain.Menu, error) {
	cacheKey := fmt.Sprintf("menus:%d", id)

	cached, err := r.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var menu domain.Menu
		if err := json.Unmarshal([]byte(cached), &menu); err == nil {
			return &menu, nil
		}
	}

	var menu domain.Menu
	if err := r.db.Preload("Restaurant").First(&menu, id).Error; err != nil {
		return nil, err
	}

	if b, err := json.Marshal(menu); err == nil {
		r.rdb.Set(ctx, cacheKey, b, 5*time.Minute)
	}

	return &menu, nil
}

func (r *repository) Create(menu *domain.Menu) error {
	if err := r.db.Create(menu).Error; err != nil {
		return err
	}
	r.invalidateListCache(menu.RestaurantID)
	return r.db.Preload("Restaurant").First(menu, menu.ID).Error
}

func (r *repository) Update(menu *domain.Menu) error {
	if err := r.db.Save(menu).Error; err != nil {
		return err
	}
	r.invalidateListCache(menu.RestaurantID)
	r.rdb.Del(ctx, fmt.Sprintf("menus:%d", menu.ID))
	return r.db.Preload("Restaurant").First(menu, menu.ID).Error
}

func (r *repository) Delete(id uint) error {
	menu, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if err := r.db.Delete(&domain.Menu{}, id).Error; err != nil {
		return err
	}
	r.invalidateListCache(menu.RestaurantID)
	r.rdb.Del(ctx, fmt.Sprintf("menus:%d", id))
	return nil
}

func (r *repository) IsOwnedBy(menuID, userID uint) bool {
	var count int64
	r.db.Model(&domain.Menu{}).
		Joins("JOIN restaurants ON restaurants.id = menus.restaurant_id").
		Where("menus.id = ? AND restaurants.created_by = ?", menuID, userID).
		Count(&count)
	return count > 0
}

func (r *repository) invalidateListCache(restaurantID uint) {
	keys, _ := r.rdb.Keys(ctx, fmt.Sprintf("menus:all:%d:*", restaurantID)).Result()
	allKeys, _ := r.rdb.Keys(ctx, "menus:all:0:*").Result()
	keys = append(keys, allKeys...)
	if len(keys) > 0 {
		r.rdb.Del(ctx, keys...)
	}
}
