package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"dailystep-backend/internal/auth"
	"dailystep-backend/internal/model"
	"dailystep-backend/internal/repository"
	"dailystep-backend/internal/transport/dto"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyUsed   = errors.New("email already used")
	ErrInvalidInput       = errors.New("invalid input")
)

type AuthService struct {
	users     *repository.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewAuthService(users *repository.UserRepository, jwtSecret string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		users:     users,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" || len(req.Password) < 8 {
		return nil, ErrInvalidInput
	}

	_, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return nil, ErrEmailAlreadyUsed
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &model.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(user.ID, s.jwtSecret, s.tokenTTL)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" || req.Password == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID, s.jwtSecret, s.tokenTTL)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	}, nil
}
