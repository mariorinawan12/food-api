package auth

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

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.usecase.Register(req)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "register success", user)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.usecase.Login(req)
	if err != nil {
		helper.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "login success", res)
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, total, err := h.usecase.GetAllUsers(page, limit)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var result []UserResponse
	for _, u := range users {
		result = append(result, UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role.Name,
			CreatedAt: u.CreatedAt,
		})
	}

	helper.Success(c, http.StatusOK, "success", gin.H{
		"data":  result,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	if err := h.usecase.ChangePassword(userID, req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "password changed successfully", nil)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	user, err := h.usecase.UpdateProfile(userID, req)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "profile updated", UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role.Name,
		CreatedAt: user.CreatedAt,
	})
}
