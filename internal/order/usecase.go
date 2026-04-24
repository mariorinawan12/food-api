package order

import (
	"errors"

	"github.com/mariorinawan12/food-api/internal/domain"
	"gorm.io/gorm"
)

type Usecase interface {
	Checkout(userID uint, req CheckoutRequest) (*domain.Order, error)
	GetMyOrders(userID uint, status string, page, limit int) ([]domain.Order, int64, error)
	GetAllOrders(status string, page, limit int) ([]domain.Order, int64, error)
	GetRestaurantOrders(adminID uint, status string, page, limit int) ([]domain.Order, int64, error)
	GetByID(userID uint, orderID uint) (*domain.Order, error)
	Pay(userID uint, orderID uint) (*domain.Order, error)
	Cancel(userID uint, orderID uint) (*domain.Order, error)
	Process(adminID uint, orderID uint) (*domain.Order, error)
	Deliver(adminID uint, orderID uint) (*domain.Order, error)
}

type usecase struct {
	repo Repository
	db   *gorm.DB
}

func NewUsecase(repo Repository, db *gorm.DB) Usecase {
	return &usecase{repo, db}
}

var validTransitions = map[string][]string{
	"payment_pending": {"paid", "cancelled"},
	"paid":            {"processing"},
	"processing":      {"delivered"},
}

func isValidTransition(current string, next string) bool {
	allowed, ok := validTransitions[current]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == next {
			return true
		}
	}
	return false
}

func (u *usecase) Checkout(userID uint, req CheckoutRequest) (*domain.Order, error) {
	var cart domain.Cart
	if err := u.db.Preload("CartItems.Menu").
		Where("id = ? AND user_id = ?", req.CartID, userID).
		First(&cart).Error; err != nil {
		return nil, errors.New("cart not found")
	}

	if len(cart.CartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	var totalPrice float64
	var orderItems []domain.OrderItem
	for _, item := range cart.CartItems {
		subtotal := item.Menu.Price * float64(item.Quantity)
		totalPrice += subtotal
		orderItems = append(orderItems, domain.OrderItem{
			MenuID:   item.MenuID,
			Quantity: item.Quantity,
			Price:    item.Menu.Price,
		})
	}

	order := &domain.Order{
		UserID:     userID,
		Address:    req.Address,
		Status:     "payment_pending",
		TotalPrice: totalPrice,
		OrderItems: orderItems,
	}

	if err := u.repo.Create(order); err != nil {
		return nil, err
	}

	u.db.Where("cart_id = ?", cart.ID).Delete(&domain.CartItem{})
	u.db.Delete(&domain.Cart{}, cart.ID)

	return u.repo.FindByID(order.ID)
}

func (u *usecase) GetMyOrders(userID uint, status string, page int, limit int) ([]domain.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return u.repo.FindAllByUserID(userID, status, page, limit)
}

func (u *usecase) GetAllOrders(status string, page int, limit int) ([]domain.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return u.repo.FindAll(status, page, limit)
}

func (u *usecase) GetRestaurantOrders(adminID uint, status string, page int, limit int) ([]domain.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return u.repo.FindAllByRestaurantAdminID(adminID, status, page, limit)
}

func (u *usecase) GetByID(userID uint, orderID uint) (*domain.Order, error) {
	order, err := u.repo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	return order, nil
}

func (u *usecase) Pay(userID uint, orderID uint) (*domain.Order, error) {
	order, err := u.repo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	if !isValidTransition(order.Status, "paid") {
		return nil, errors.New("cannot pay order with status " + order.Status)
	}

	if err := u.repo.UpdateStatus(orderID, "paid"); err != nil {
		return nil, err
	}
	return u.repo.FindByID(orderID)
}

func (u *usecase) Cancel(userID uint, orderID uint) (*domain.Order, error) {
	order, err := u.repo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	if !isValidTransition(order.Status, "cancelled") {
		return nil, errors.New("cannot cancel order with status " + order.Status)
	}

	if err := u.repo.UpdateStatus(orderID, "cancelled"); err != nil {
		return nil, err
	}
	return u.repo.FindByID(orderID)
}

func (u *usecase) Process(adminID uint, orderID uint) (*domain.Order, error) {
	order, err := u.repo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if !u.repo.IsOrderOwnedByAdmin(adminID, orderID) {
		return nil, errors.New("unauthorized, order does not belong to your restaurant")
	}

	if !isValidTransition(order.Status, "processing") {
		return nil, errors.New("cannot process order with status " + order.Status)
	}

	if err := u.repo.UpdateStatus(orderID, "processing"); err != nil {
		return nil, err
	}
	return u.repo.FindByID(orderID)
}

func (u *usecase) Deliver(adminID uint, orderID uint) (*domain.Order, error) {
	order, err := u.repo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if !u.repo.IsOrderOwnedByAdmin(adminID, orderID) {
		return nil, errors.New("unauthorized, order does not belong to your restaurant")
	}

	if !isValidTransition(order.Status, "delivered") {
		return nil, errors.New("cannot deliver order with status " + order.Status)
	}

	if err := u.repo.UpdateStatus(orderID, "delivered"); err != nil {
		return nil, err
	}
	return u.repo.FindByID(orderID)
}
