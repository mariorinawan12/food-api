package review

import (
	"errors"

	"github.com/mariorinawan12/food-api/internal/domain"
)

type Usecase interface {
	GetAllByRestaurant(restaurantID uint, page, limit int) ([]domain.Review, int64, error)
	Create(userID uint, req CreateRequest) (*domain.Review, error)
	Update(userID, reviewID uint, req UpdateRequest) (*domain.Review, error)
	Delete(userID, reviewID uint) error
}

type usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) Usecase {
	return &usecase{repo}
}

func (u *usecase) GetAllByRestaurant(restaurantID uint, page, limit int) ([]domain.Review, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return u.repo.FindAllByRestaurantID(restaurantID, page, limit)
}

func (u *usecase) Create(userID uint, req CreateRequest) (*domain.Review, error) {
	if !u.repo.HasDeliveredOrder(userID, req.RestaurantID) {
		return nil, errors.New("you can only review a restaurant after a delivered order")
	}

	existing, err := u.repo.FindByUserAndRestaurant(userID, req.RestaurantID)
	if err == nil && existing.ID != 0 {
		return nil, errors.New("you have already reviewed this restaurant")
	}

	review := &domain.Review{
		UserID:       userID,
		RestaurantID: req.RestaurantID,
		Rating:       req.Rating,
		Comment:      req.Comment,
	}

	if err := u.repo.Create(review); err != nil {
		return nil, err
	}

	return u.repo.FindByID(review.ID)
}

func (u *usecase) Update(userID, reviewID uint, req UpdateRequest) (*domain.Review, error) {
	review, err := u.repo.FindByID(reviewID)
	if err != nil {
		return nil, errors.New("review not found")
	}

	if review.UserID != userID {
		return nil, errors.New("unauthorized, you are not the owner")
	}

	review.Rating = req.Rating
	review.Comment = req.Comment

	if err := u.repo.Update(review); err != nil {
		return nil, err
	}

	return u.repo.FindByID(reviewID)
}

func (u *usecase) Delete(userID, reviewID uint) error {
	review, err := u.repo.FindByID(reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	if review.UserID != userID {
		return errors.New("unauthorized, you are not the owner")
	}

	return u.repo.Delete(reviewID)
}
