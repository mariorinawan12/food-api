package auth

import (
	"errors"

	"github.com/mariorinawan12/food-api/internal/domain"
	"github.com/mariorinawan12/food-api/internal/helper"
)

type Usecase interface {
	Register(req RegisterRequest) (*domain.User, error)
	Login(req LoginRequest) (*LoginResponse, error)
	GetAllUsers(page int, limit int) ([]domain.User, int64, error)
	ChangePassword(userID uint, req ChangePasswordRequest) error
	UpdateProfile(userID uint, req UpdateProfileRequest) (*domain.User, error)
}

type usecase struct {
	repo Repository
}

func NewUseCase(repo Repository) Usecase {
	return &usecase{repo}
}

func (u *usecase) Register(req RegisterRequest) (*domain.User, error) {
	role, err := u.repo.FindRoleByName(req.Role)
	if err != nil {
		return nil, errors.New("default role not found")
	}

	hashed, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
		RoleID:   role.ID,
	}

	if err := u.repo.Create(user); err != nil {
		return nil, errors.New("email already exists")
	}

	return user, nil
}

func (u *usecase) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := u.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !helper.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := helper.GenerateToken(user.ID, user.Role.Name)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: token, User: *user}, nil
}

func (u *usecase) GetAllUsers(page, limit int) ([]domain.User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return u.repo.FindAll(page, limit)
}

func (u *usecase) ChangePassword(userID uint, req ChangePasswordRequest) error {
	user, err := u.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !helper.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	hashed, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return u.repo.UpdatePassword(userID, hashed)
}

func (u *usecase) UpdateProfile(userID uint, req UpdateProfileRequest) (*domain.User, error) {
	user, err := u.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !helper.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("incorrect password")
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" && req.Email != user.Email {
		existing, err := u.repo.FindByEmail(req.Email)
		if err == nil && existing.ID != userID {
			return nil, errors.New("email already in use")
		}
		user.Email = req.Email
	}

	if err := u.repo.UpdateProfile(userID, user.Name, user.Email); err != nil {
		return nil, err
	}

	return u.repo.FindByID(userID)
}
