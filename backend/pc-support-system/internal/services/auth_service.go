package services

import (
	"context"
	"errors"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/validation"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  repo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {

	var err error

	// Validate & normalize input
	req.FirstName, err = validation.ValidateName(req.FirstName, "first name")
	if err != nil {
		return err
	}

	if req.LastName != "" {
		req.LastName, err = validation.ValidateName(req.LastName, "last name")
		if err != nil {
			return err
		}
	}

	req.Phone, err = validation.ValidatePhone(req.Phone)
	if err != nil {
		return err
	}

	req.Email, err = validation.ValidateEmail(req.Email, true)
	if err != nil {
		return err
	}

	// Password validation
	if err := validation.ValidatePassword(req.Password); err != nil {
		return err
	}

	// Check if phone already exists
	existingUser, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.New("phone number already exists")
	}

	// Check if email already exists (only if provided)
	if req.Email != "" {

		existingUser, err = s.userRepo.FindByEmail(ctx, req.Email)
		if err != nil {
			return err
		}

		if existingUser != nil {
			return errors.New("email already exists")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	now := time.Now()

	// Create user
	user := &models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hashedPassword),

		Role:  models.RoleCustomer,
		State: models.UserActive,

		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save
	return s.userRepo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*models.User, error) {

	var err error

	req.Email, err = validation.ValidateEmail(req.Email, true)
	if err != nil {
		return nil, err
	}

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

func (s *AuthService) UpdateRefreshToken(
	ctx context.Context,
	userID bson.ObjectID,
	token string,
	expiresAt time.Time,
) error {
	return s.userRepo.UpdateRefreshToken(
		ctx,
		userID,
		token,
		expiresAt,
	)
}

func GenerateRefreshToken(
	user *models.User,
	secret string,
) (string, time.Time, error) {

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"type":    "refresh",
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (s *AuthService) ClearRefreshToken(
	ctx context.Context,
	userID bson.ObjectID,
) error {

	return s.userRepo.ClearRefreshToken(
		ctx,
		userID,
	)
}

func (s *AuthService) FindByID(
	ctx context.Context,
	userID bson.ObjectID,
) (*models.User, error) {

	return s.userRepo.FindByID(ctx, userID)
}
