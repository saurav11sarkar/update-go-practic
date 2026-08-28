package user

import (
	"errors"
	"go-prictic/internal/auth"
	"go-prictic/internal/user/dto"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service interface {
	CreateUser(req dto.UserRequest) (*dto.UserResponse, error)
	LoginUser(req dto.UserLoginRequest) (*dto.UserLoginResponse, error)
	GetUserByID(id string) (*dto.UserResponse, error)
}

type service struct {
	accessTokenService  auth.JwtService
	refreshTokenService auth.JwtService
	repo                Repository
}

func NewService(repo Repository, accessTokenService, refreshTokenService auth.JwtService) Service {
	return &service{
		repo:                repo,
		accessTokenService:  accessTokenService,
		refreshTokenService: refreshTokenService,
	}
}

func (s *service) CreateUser(req dto.UserRequest) (*dto.UserResponse, error) {
	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		// Role:     Role(req.Role),
	}
	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     string(user.Role),
		IsActive: user.IsActive,
		CreateAt: user.CreateAt,
		UpdateAt: user.UpdateAt,
	}, nil
}

func (s *service) LoginUser(req dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	accessToken, err := s.accessTokenService.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.refreshTokenService.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}
	return &dto.UserLoginResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Role:         string(user.Role),
		IsActive:     user.IsActive,
		CreateAt:     user.CreateAt,
		UpdateAt:     user.UpdateAt,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) GetUserByID(id string) (*dto.UserResponse, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     string(user.Role),
		IsActive: user.IsActive,
		CreateAt: user.CreateAt,
		UpdateAt: user.UpdateAt,
	}, nil
}
