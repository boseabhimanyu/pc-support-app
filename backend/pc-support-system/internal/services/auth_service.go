package services

import (
	"context"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: repo,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,

		Role:   models.RoleCustomer,
		Active: true,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.userRepo.Create(ctx, user)
}

// Login()

// HashPassword()

// ComparePassword()

// GenerateJWT()

// ValidateRole()
