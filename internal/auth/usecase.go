package auth

import (
	"errors"

	"github.com/mariorinawan12/food-api/internal/domain"
	"github.com/mariorinawan12/food-api/internal/helper"
)

type Usecase interface {
	Register(req RegisterRequest) (*domain.User, error)
	Login(req LoginRequest) (*LoginResponse, error)
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
