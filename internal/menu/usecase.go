package menu

import (
	"errors"

	"github.com/mariorinawan12/food-api/internal/domain"
)

type RestaurantRepository interface {
	IsOwnedBy(restaurantID, userID uint) bool
}

type Usecase interface {
	GetAll(restaurantID uint, search string, page, limit int) ([]domain.Menu, int64, error)
	GetByID(id uint) (*domain.Menu, error)
	Create(req CreateRequest, userID uint, role string) (*domain.Menu, error)
	Update(id uint, req UpdateRequest, userID uint, role string) (*domain.Menu, error)
	Delete(id uint, userID uint, role string) error
}

type usecase struct {
	repo           Repository
	restaurantRepo RestaurantRepository
}

func NewUsecase(repo Repository, restaurantRepo RestaurantRepository) Usecase {
	return &usecase{repo, restaurantRepo}
}

func (u *usecase) GetAll(restaurantID uint, search string, page, limit int) ([]domain.Menu, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return u.repo.FindAll(restaurantID, search, page, limit)
}

func (u *usecase) GetByID(id uint) (*domain.Menu, error) {
	menu, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("menu not found")
	}
	return menu, nil
}

func (u *usecase) Create(req CreateRequest, userID uint, role string) (*domain.Menu, error) {

	if role != "super_admin" {
		if !u.restaurantRepo.IsOwnedBy(req.RestaurantID, userID) {
			return nil, errors.New("unauthorized, restaurant does not belong to you")
		}
	}

	menu := &domain.Menu{
		RestaurantID: req.RestaurantID,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		Category:     req.Category,
	}
	if err := u.repo.Create(menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (u *usecase) Update(id uint, req UpdateRequest, userID uint, role string) (*domain.Menu, error) {
	menu, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("menu not found")
	}

	if role != "super_admin" {
		if !u.repo.IsOwnedBy(id, userID) {
			return nil, errors.New("unauthorized, you are not the owner")
		}
	}

	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Description != "" {
		menu.Description = req.Description
	}
	if req.Price > 0 {
		menu.Price = req.Price
	}
	if req.Category != "" {
		menu.Category = req.Category
	}

	if err := u.repo.Update(menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (u *usecase) Delete(id uint, userID uint, role string) error {
	if _, err := u.repo.FindByID(id); err != nil {
		return errors.New("menu not found")
	}

	if role != "super_admin" {
		if !u.repo.IsOwnedBy(id, userID) {
			return errors.New("unauthorized, you are not the owner")
		}
	}

	return u.repo.Delete(id)
}
