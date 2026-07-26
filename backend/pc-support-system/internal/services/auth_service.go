package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
	"golang.org/x/crypto/bcrypt"
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

	// 1. Check if phone already exists
	existingUser, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.New("phone number already exists")
	}

	// 2. Check if email already exists (only if email is provided)
	existingUser, err = s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	// 3. Create the user
	user := &models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleCustomer,
		State:        models.UserActive,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 4. Save to MongoDB
	return s.userRepo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*models.User, error) {

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if account is active
	if user.State != models.UserActive {
		return nil, errors.New("account is inactive")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
