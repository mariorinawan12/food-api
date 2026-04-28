package cart

import (
	"errors"

	"github.com/mariorinawan12/food-api/internal/domain"
)

type MenuRepository interface {
	FindByID(id uint) (*domain.Menu, error)
}

type Usecase interface {
	GetAllCarts(userID uint, page int, limit int) ([]domain.Cart, int64, error)
	GetCartByID(userID uint, cartID uint) (*domain.Cart, error)
	CreateCart(userID uint, req CreateCartRequest) (*domain.Cart, error)
	DeleteCart(userID uint, cartID uint) error
	AddItem(userID uint, cartID uint, req AddItemRequest) (*domain.Cart, error)
	UpdateItem(userID uint, cartID uint, itemID uint, req UpdateItemRequest) (*domain.Cart, error)
	DeleteItem(userID uint, cartID uint, itemID uint) (*domain.Cart, error)
}

type usecase struct {
	repo     Repository
	menuRepo MenuRepository
}

func NewUsecase(repo Repository, menuRepo MenuRepository) Usecase {
	return &usecase{repo, menuRepo}
}

func (u *usecase) GetAllCarts(userID uint, page, limit int) ([]domain.Cart, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return u.repo.FindAllByUserID(userID, page, limit)
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

func (u *usecase) DeleteCart(userID uint, cartID uint) error {
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

	// ← validasi menu yang mau ditambah
	menu, err := u.menuRepo.FindByID(req.MenuID)
	if err != nil {
		return nil, errors.New("menu not found")
	}
	if menu.RestaurantID != cart.RestaurantID {
		return nil, errors.New("menu does not not belong to this restaurant")
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

func (u *usecase) UpdateItem(userID uint, cartID uint, itemID uint, req UpdateItemRequest) (*domain.Cart, error) {
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

func (u *usecase) DeleteItem(userID uint, cartID uint, itemID uint) (*domain.Cart, error) {
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
