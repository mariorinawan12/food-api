package restaurant

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
	categoryID, _ := strconv.Atoi(c.Query("category_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	restaurants, total, err := h.usecase.GetAll(search, uint(categoryID), page, limit)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var result []RestaurantResponse
	for _, r := range restaurants {
		result = append(result, RestaurantResponse{
			ID:           r.ID,
			Name:         r.Name,
			Description:  r.Description,
			Address:      r.Address,
			CategoryName: r.Category.Name,
			CreatedBy:    r.CreatedBy,
			CreatedAt:    r.CreatedAt,
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	r, err := h.usecase.GetByID(uint(id))
	if err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "success", RestaurantResponse{
		ID:           r.ID,
		Name:         r.Name,
		Description:  r.Description,
		Address:      r.Address,
		CategoryName: r.Category.Name,
		CreatedBy:    r.CreatedBy,
		CreatedAt:    r.CreatedAt,
	})
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	r, err := h.usecase.Create(req, userID)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "restaurant created", RestaurantResponse{
		ID:           r.ID,
		Name:         r.Name,
		Description:  r.Description,
		Address:      r.Address,
		CategoryName: r.Category.Name,
		CreatedBy:    r.CreatedBy,
		CreatedAt:    r.CreatedAt,
	})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
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
	r, err := h.usecase.Update(uint(id), req, userID)
	if err != nil {
		helper.Error(c, http.StatusForbidden, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "restaurant updated", RestaurantResponse{
		ID:           r.ID,
		Name:         r.Name,
		Description:  r.Description,
		Address:      r.Address,
		CategoryName: r.Category.Name,
		CreatedBy:    r.CreatedBy,
		CreatedAt:    r.CreatedAt,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	userID := c.GetUint("user_id")
	if err := h.usecase.Delete(uint(id), userID); err != nil {
		helper.Error(c, http.StatusForbidden, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "restaurant deleted", nil)
}

func (h *Handler) GetMyRestaurants(c *gin.Context) {
	userID := c.GetUint("user_id")

	restaurants, err := h.usecase.GetMyRestaurants(userID)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var result []RestaurantResponse
	for _, r := range restaurants {
		result = append(result, RestaurantResponse{
			ID:           r.ID,
			Name:         r.Name,
			Description:  r.Description,
			Address:      r.Address,
			CategoryName: r.Category.Name,
			CreatedBy:    r.CreatedBy,
			CreatedAt:    r.CreatedAt,
		})
	}

	helper.Success(c, http.StatusOK, "success", result)
}
