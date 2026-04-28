package cart

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mariorinawan12/food-api/internal/domain"
	"github.com/mariorinawan12/food-api/internal/helper"
)

type Handler struct {
	usecase  Usecase
	validate *validator.Validate
}

func NewHandler(usecase Usecase) *Handler {
	return &Handler{usecase, validator.New()}
}

func buildCartResponse(cart *domain.Cart) CartResponse {
	var items []CartItemResponse
	var totalPrice float64

	for _, item := range cart.CartItems {
		subtotal := item.Menu.Price * float64(item.Quantity)
		totalPrice += subtotal
		items = append(items, CartItemResponse{
			ID:       item.ID,
			MenuID:   item.MenuID,
			MenuName: item.Menu.Name,
			Price:    item.Menu.Price,
			Quantity: item.Quantity,
			Subtotal: subtotal,
		})
	}

	return CartResponse{
		ID:             cart.ID,
		UserID:         cart.UserID,
		RestaurantID:   cart.RestaurantID,
		RestaurantName: cart.Restaurant.Name,
		Items:          items,
		TotalPrice:     totalPrice,
		CreatedAt:      cart.CreatedAt,
	}
}

func (h *Handler) GetAllCarts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userID := c.GetUint("user_id")

	carts, total, err := h.usecase.GetAllCarts(userID, page, limit)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var result []CartResponse
	for _, cart := range carts {
		result = append(result, buildCartResponse(&cart))
	}

	helper.Success(c, http.StatusOK, "success", gin.H{
		"data":  result,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) GetCartByID(c *gin.Context) {
	cartID, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	userID := c.GetUint("user_id")
	cart, err := h.usecase.GetCartByID(userID, uint(cartID))
	if err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "success", buildCartResponse(cart))
}

func (h *Handler) CreateCart(c *gin.Context) {
	restaurantID, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid restaurant id")
		return
	}

	userID := c.GetUint("user_id")
	cart, err := h.usecase.CreateCart(userID, CreateCartRequest{
		RestaurantID: uint(restaurantID),
	})
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "cart created", buildCartResponse(cart))
}

func (h *Handler) DeleteCart(c *gin.Context) {
	cartID, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	userID := c.GetUint("user_id")
	if err := h.usecase.DeleteCart(userID, uint(cartID)); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "cart deleted", nil)
}

func (h *Handler) AddItem(c *gin.Context) {
	cartID, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	cart, err := h.usecase.AddItem(userID, uint(cartID), req)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "item added", buildCartResponse(cart))
}

func (h *Handler) UpdateItem(c *gin.Context) {
	cartID, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	itemID, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid item id")
		return
	}

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	cart, err := h.usecase.UpdateItem(userID, uint(cartID), uint(itemID), req)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "item updated", buildCartResponse(cart))
}

func (h *Handler) DeleteItem(c *gin.Context) {
	cartID, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid cart id")
		return
	}

	itemID, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid item id")
		return
	}

	userID := c.GetUint("user_id")
	cart, err := h.usecase.DeleteItem(userID, uint(cartID), uint(itemID))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "item deleted", buildCartResponse(cart))
}
