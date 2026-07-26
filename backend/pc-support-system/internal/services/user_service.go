package services

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetProfile(
	ctx context.Context,
	id bson.ObjectID,
) (*dto.UserResponse, error) {

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	resp := dto.ToUserResponse(user)

	return &resp, nil
}

func (s *UserService) UpdateProfile(
	ctx context.Context,
	id bson.ObjectID,
	req dto.UpdateProfileRequest,
) (*dto.UserResponse, error) {

	// Find current user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	// Check email uniqueness
	if req.Email != nil {

		if strings.TrimSpace(*req.Email) == "" {
			return nil, errors.New("email cannot be empty")
		}

		if *req.Email != user.Email {
			existing, err := s.userRepo.FindByEmail(ctx, *req.Email)
			if err != nil {
				return nil, err
			}

			if existing != nil {
				return nil, errors.New("email already in use with another user")
			}
		}

		user.Email = *req.Email
	}

	// Check phone uniqueness
	if req.Phone != nil {

		if strings.TrimSpace(*req.Phone) == "" {
			return nil, errors.New("phone cannot be empty")
		}

		if *req.Phone != user.Phone {
			existing, err := s.userRepo.FindByPhone(ctx, *req.Phone)
			if err != nil {
				return nil, err
			}

			if existing != nil {
				return nil, errors.New("phone number already in use with another user")
			}
		}

		user.Phone = *req.Phone
	}

	user.UpdatedAt = time.Now()

	// Save
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	resp := dto.ToUserResponse(user)

	return &resp, nil
}

func (s *UserService) SearchCustomers(
	ctx context.Context,
	query string,
) ([]dto.CustomerSearchResponse, error) {

	query = strings.TrimSpace(query)

	if query == "" {
		return nil, errors.New("search query is required")
	}

	users, err := s.userRepo.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	return dto.ToCustomerSearchResponses(users), nil
}

func (s *UserService) CreateCustomer(
	ctx context.Context,
	createdBy string,
	req dto.CreateCustomerRequest,
) (*dto.CustomerResponse, error) {

	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Phone = strings.TrimSpace(req.Phone)
	if matched, _ := regexp.MatchString(`^[0-9]+$`, req.Phone); !matched {
		return nil, errors.New("invalid phone number")
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// Phone must be unique
	existingPhone, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}

	if existingPhone != nil {
		return nil, errors.New("phone number already exists")
	}

	// Email is optional but must be unique if provided
	if req.Email != "" {

		existingEmail, err := s.userRepo.FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}

		if existingEmail != nil {
			return nil, errors.New("email already exists")
		}
	}

	createdByID, err := bson.ObjectIDFromHex(createdBy)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	// Verify creator exists
	creator, err := s.userRepo.FindByID(ctx, createdByID)
	if err != nil {
		return nil, err
	}

	if creator == nil {
		return nil, errors.New("creator not found")
	}

	now := time.Now()

	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Email:     req.Email,

		Role:  models.RoleCustomer,
		State: models.UserActive,

		PasswordHash: "",
		CreatedByID:  &createdByID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	resp := dto.ToCustomerResponse(user)

	return &resp, nil
}

func (s *UserService) SetCustomerPassword(
	ctx context.Context,
	customerID string,
	req dto.SetCustomerPasswordRequest,
) error {

	id, err := bson.ObjectIDFromHex(customerID)
	if err != nil {
		return errors.New("invalid customer id")
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("customer not found")
	}

	if user.Role != models.RoleCustomer {
		return errors.New("invalid customer")
	}

	if strings.TrimSpace(req.Password) == "" {
		return errors.New("password is required")
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(
		ctx,
		id,
		string(hash),
	)
}
