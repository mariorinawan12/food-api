package review

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

func buildReviewResponse(r *domain.Review) ReviewResponse {
	return ReviewResponse{
		ID:           r.ID,
		UserID:       r.UserID,
		UserName:     r.User.Name,
		RestaurantID: r.RestaurantID,
		Rating:       r.Rating,
		Comment:      r.Comment,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func (h *Handler) GetAllByRestaurant(c *gin.Context) {
	restaurantID, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid restaurant id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, total, err := h.usecase.GetAllByRestaurant(uint(restaurantID), page, limit)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var result []ReviewResponse
	for _, r := range reviews {
		result = append(result, buildReviewResponse(&r))
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
		helper.Error(c, http.StatusBadRequest, "invalid restaurant ID")
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
	review, err := h.usecase.Create(userID, req)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "review created", buildReviewResponse(review))
}

func (h *Handler) Update(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("review_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid review id")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	review, err := h.usecase.Update(userID, uint(reviewID), req)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "review updated", buildReviewResponse(review))
}

func (h *Handler) Delete(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("review_id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	userID := c.GetUint("user_id")
	if err := h.usecase.Delete(userID, uint(reviewID)); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "review deleted", nil)
}
