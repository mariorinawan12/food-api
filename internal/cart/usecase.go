package cart

import (
	"errors"

	"github.com/mariorinawan12/food-api/internal/domain"
)

type Usecase interface {
	GetAllCarts(userID uint) ([]domain.Cart, error)
	GetCartByID(userID uint, cartID uint) (*domain.Cart, error)
	CreateCart(userID uint, req CreateCartRequest) (*domain.Cart, error)
	DeleteCart(userID uint, cartID uint) error
	AddItem(userID uint, cartID uint, req AddItemRequest) (*domain.Cart, error)
	UpdateItem(userID uint, cartID uint, itemID uint, req UpdateItemRequest) (*domain.Cart, error)
	DeleteItem(userID uint, cartID uint, itemID uint) (*domain.Cart, error)
}

type usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) Usecase {
	return &usecase{repo}
}

func (u *usecase) GetAllCarts(userID uint) ([]domain.Cart, error) {
	return u.repo.FindAllByUserID(userID)
}

func (u *usecase) GetCartByID(userID uint, cartID uint) (*domain.Cart, error) {
	cart, err := u.repo.FindByID(cartID)
	if err != nil {
		return nil, errors.New("cart not found")
	}
	if cart.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	return cart, nil
}

func (u *usecase) CreateCart(userID uint, req CreateCartRequest) (*domain.Cart, error) {
	// cek apakah sudah ada cart untuk restoran ini
	existing, err := u.repo.FindByUserIDAndRestaurantID(userID, req.RestaurantID)
	if err == nil && existing.ID != 0 {
		return nil, errors.New("cart for this restaurant already exists")
	}

	cart := &domain.Cart{
		UserID:       userID,
		RestaurantID: req.RestaurantID,
	}

	if err := u.repo.Create(cart); err != nil {
		return nil, err
	}
	return cart, nil
}

func (u *usecase) DeleteCart(userID, cartID uint) error {
	cart, err := u.repo.FindByID(cartID)
	if err != nil {
		return errors.New("cart not found")
	}
	if cart.UserID != userID {
		return errors.New("unauthorized")
	}
	return u.repo.Delete(cartID)
}

func (u *usecase) AddItem(userID, cartID uint, req AddItemRequest) (*domain.Cart, error) {
	cart, err := u.repo.FindByID(cartID)
	if err != nil {
		return nil, errors.New("cart not found")
	}
	if cart.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// validasi menu dari restoran yang sama
	for _, item := range cart.CartItems {
		if item.Menu.RestaurantID != cart.RestaurantID {
			return nil, errors.New("menu does not belong to this restaurant")
		}
	}

	// cek apakah menu sudah ada di cart
	existing, err := u.repo.FindCartItemByMenuID(cartID, req.MenuID)
	if err == nil {
		existing.Quantity += req.Quantity
		if err := u.repo.UpdateItem(existing); err != nil {
			return nil, err
		}
	} else {
		item := &domain.CartItem{
			CartID:   cartID,
			MenuID:   req.MenuID,
			Quantity: req.Quantity,
		}
		if err := u.repo.AddItem(item); err != nil {
			return nil, err
		}
	}

	return u.repo.FindByID(cartID)
}

func (u *usecase) UpdateItem(userID, cartID, itemID uint, req UpdateItemRequest) (*domain.Cart, error) {
	cart, err := u.repo.FindByID(cartID)
	if err != nil {
		return nil, errors.New("cart not found")
	}
	if cart.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	item, err := u.repo.FindCartItemByID(itemID)
	if err != nil {
		return nil, errors.New("item not found")
	}
	if item.CartID != cartID {
		return nil, errors.New("item does not belong to this cart")
	}

	item.Quantity = req.Quantity
	if err := u.repo.UpdateItem(item); err != nil {
		return nil, err
	}

	return u.repo.FindByID(cartID)
}

func (u *usecase) DeleteItem(userID, cartID, itemID uint) (*domain.Cart, error) {
	cart, err := u.repo.FindByID(cartID)
	if err != nil {
		return nil, errors.New("cart not found")
	}
	if cart.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	item, err := u.repo.FindCartItemByID(itemID)
	if err != nil {
		return nil, errors.New("item not found")
	}
	if item.CartID != cartID {
		return nil, errors.New("item does not belong to this cart")
	}

	if err := u.repo.DeleteItem(itemID); err != nil {
		return nil, err
	}

	return u.repo.FindByID(cartID)
}
