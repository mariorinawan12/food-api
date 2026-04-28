package restaurant

import (
	"errors"

	"github.com/mariorinawan12/food-api/internal/domain"
)

type Usecase interface {
	GetAll(search string, categoryID uint, page, limit int) ([]domain.Restaurant, int64, error)
	GetByID(id uint) (*domain.Restaurant, error)
	Create(req CreateRequest, userID uint) (*domain.Restaurant, error)
	Update(id uint, req UpdateRequest, userID uint, role string) (*domain.Restaurant, error)
	Delete(id uint, userID uint, role string) error
	GetMyRestaurants(userID uint) ([]domain.Restaurant, error)
}

type usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) Usecase {
	return &usecase{repo}
}

func (u *usecase) GetAll(search string, categoryID uint, page, limit int) ([]domain.Restaurant, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return u.repo.FindAll(search, categoryID, page, limit)
}

func (u *usecase) GetByID(id uint) (*domain.Restaurant, error) {
	restaurant, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}
	return restaurant, nil
}

func (u *usecase) Create(req CreateRequest, userID uint) (*domain.Restaurant, error) {
	restaurant := &domain.Restaurant{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		CategoryID:  req.CategoryID,
		CreatedBy:   userID,
	}
	if err := u.repo.Create(restaurant); err != nil {
		return nil, err
	}
	return restaurant, nil
}

func (u *usecase) Update(id uint, req UpdateRequest, userID uint, role string) (*domain.Restaurant, error) {
	restaurant, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	if role != "super_admin" {
		if !u.repo.IsOwnedBy(id, userID) {
			return nil, errors.New("unauthorized, you are not the owner")
		}
	}
	if req.Name != "" {
		restaurant.Name = req.Name
	}
	if req.Description != "" {
		restaurant.Description = req.Description
	}
	if req.Address != "" {
		restaurant.Address = req.Address
	}
	if req.CategoryID != 0 {
		restaurant.CategoryID = req.CategoryID
	}

	if err := u.repo.Update(restaurant); err != nil {
		return nil, err
	}
	return restaurant, nil
}

func (u *usecase) Delete(id uint, userID uint, role string) error {
	if _, err := u.repo.FindByID(id); err != nil {
		return errors.New("restaurant not found")
	}

	if role != "super_admin" {
		if !u.repo.IsOwnedBy(id, userID) {
			return errors.New("unauthorized, you are not the owner")
		}
	}

	return u.repo.Delete(id)
}

func (u *usecase) GetMyRestaurants(userID uint) ([]domain.Restaurant, error) {
	return u.repo.FindByUserID(userID)
}
