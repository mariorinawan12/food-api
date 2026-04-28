package menu

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mariorinawan12/food-api/internal/helper"
)

type Handler struct {
	usecase  Usecase
	validate *validator.Validate
}

func NewHandler(usecase Usecase) *Handler {
	return &Handler{usecase, validator.New()}
}

func (h *Handler) GetAll(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	restaurantID, _ := strconv.Atoi(c.Query("restaurant_id"))

	menus, total, err := h.usecase.GetAll(uint(restaurantID), search, page, limit)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var result []MenuResponse
	for _, m := range menus {
		result = append(result, MenuResponse{
			ID:             m.ID,
			RestaurantID:   m.RestaurantID,
			RestaurantName: m.Restaurant.Name,
			Name:           m.Name,
			Description:    m.Description,
			Price:          m.Price,
			Category:       m.Category,
			CreatedAt:      m.CreatedAt,
		})
	}

	helper.Success(c, http.StatusOK, "success", gin.H{
		"data":  result,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	menuID, err := strconv.Atoi(c.Param("menu_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	m, err := h.usecase.GetByID(uint(menuID))
	if err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "success", MenuResponse{
		ID:             m.ID,
		RestaurantID:   m.RestaurantID,
		RestaurantName: m.Restaurant.Name,
		Name:           m.Name,
		Description:    m.Description,
		Price:          m.Price,
		Category:       m.Category,
		CreatedAt:      m.CreatedAt,
	})
}

func (h *Handler) GetByRestaurant(c *gin.Context) {
	restaurantID, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid restaurant id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	menus, total, err := h.usecase.GetAll(uint(restaurantID), search, page, limit)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var result []MenuResponse
	for _, m := range menus {
		result = append(result, MenuResponse{
			ID:             m.ID,
			RestaurantID:   m.RestaurantID,
			RestaurantName: m.Restaurant.Name,
			Name:           m.Name,
			Description:    m.Description,
			Price:          m.Price,
			Category:       m.Category,
			CreatedAt:      m.CreatedAt,
		})
	}

	helper.Success(c, http.StatusOK, "success", gin.H{
		"data":  result,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) Create(c *gin.Context) {

	restaurantID, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid restaurant id")
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	req.RestaurantID = uint(restaurantID)

	userID := c.GetUint("user_id")
	role := c.GetString("role")
	m, err := h.usecase.Create(req, userID, role)
	if err != nil {
		helper.Error(c, http.StatusForbidden, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "menu created", MenuResponse{
		ID:             m.ID,
		RestaurantID:   m.RestaurantID,
		RestaurantName: m.Restaurant.Name,
		Name:           m.Name,
		Description:    m.Description,
		Price:          m.Price,
		Category:       m.Category,
		CreatedAt:      m.CreatedAt,
	})
}

func (h *Handler) Update(c *gin.Context) {
	menuID, err := strconv.Atoi(c.Param("menu_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")
	m, err := h.usecase.Update(uint(menuID), req, userID, role)
	if err != nil {
		helper.Error(c, http.StatusForbidden, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "menu updated", MenuResponse{
		ID:             m.ID,
		RestaurantID:   m.RestaurantID,
		RestaurantName: m.Restaurant.Name,
		Name:           m.Name,
		Description:    m.Description,
		Price:          m.Price,
		Category:       m.Category,
		CreatedAt:      m.CreatedAt,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	menuID, err := strconv.Atoi(c.Param("menu_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")
	if err := h.usecase.Delete(uint(menuID), userID, role); err != nil {
		helper.Error(c, http.StatusForbidden, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "menu deleted", nil)
}
