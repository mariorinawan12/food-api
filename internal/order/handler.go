package order

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

func buildOrderResponse(order *domain.Order) OrderResponse {
	var items []OrderItemResponse
	for _, item := range order.OrderItems {
		items = append(items, OrderItemResponse{
			ID:       item.ID,
			MenuID:   item.MenuID,
			MenuName: item.Menu.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
			Subtotal: item.Price * float64(item.Quantity),
		})
	}
	return OrderResponse{
		ID:         order.ID,
		UserID:     order.UserID,
		Address:    order.Address,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
		Items:      items,
		CreatedAt:  order.CreatedAt,
	}
}

func (h *Handler) Checkout(c *gin.Context) {
	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	order, err := h.usecase.Checkout(userID, req)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "checkout success", buildOrderResponse(order))
}

func (h *Handler) GetMyOrders(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userID := c.GetUint("user_id")

	orders, total, err := h.usecase.GetMyOrders(userID, status, page, limit)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var result []OrderResponse
	for _, o := range orders {
		result = append(result, buildOrderResponse(&o))
	}

	helper.Success(c, http.StatusOK, "success", gin.H{
		"data":  result,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) GetAllOrders(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, total, err := h.usecase.GetAllOrders(status, page, limit)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var result []OrderResponse
	for _, o := range orders {
		result = append(result, buildOrderResponse(&o))
	}

	helper.Success(c, http.StatusOK, "success", gin.H{
		"data":  result,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) GetRestaurantOrders(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	adminID := c.GetUint("user_id")

	orders, total, err := h.usecase.GetRestaurantOrders(adminID, status, page, limit)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var result []OrderResponse
	for _, o := range orders {
		result = append(result, buildOrderResponse(&o))
	}

	helper.Success(c, http.StatusOK, "success", gin.H{
		"data":  result,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	userID := c.GetUint("user_id")
	order, err := h.usecase.GetByID(userID, uint(id))
	if err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "success", buildOrderResponse(order))
}

func (h *Handler) Pay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	userID := c.GetUint("user_id")
	order, err := h.usecase.Pay(userID, uint(id))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "payment success", buildOrderResponse(order))
}

func (h *Handler) Cancel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	userID := c.GetUint("user_id")
	order, err := h.usecase.Cancel(userID, uint(id))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "order cancelled", buildOrderResponse(order))
}

func (h *Handler) Process(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	adminID := c.GetUint("user_id")
	order, err := h.usecase.Process(adminID, uint(id))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "order processing", buildOrderResponse(order))
}

func (h *Handler) Deliver(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	adminID := c.GetUint("user_id")
	order, err := h.usecase.Deliver(adminID, uint(id))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "order delivered", buildOrderResponse(order))
}
