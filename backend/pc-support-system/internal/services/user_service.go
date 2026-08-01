package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/validation"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo     repository.UserRepository
	auditService *AuditService
}

func NewUserService(
	userRepo repository.UserRepository,
	auditService *AuditService,
) *UserService {
	return &UserService{
		userRepo:     userRepo,
		auditService: auditService,
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

	// First Name
	if req.FirstName != nil {

		firstName, err := validation.ValidateName(*req.FirstName)
		if err != nil {
			return nil, err
		}

		user.FirstName = firstName
	}

	// Last Name
	if req.LastName != nil {

		lastName, err := validation.ValidateName(*req.LastName)
		if err != nil {
			return nil, err
		}

		user.LastName = lastName
	}

	// Email
	if req.Email != nil {

		email, err := validation.ValidateEmail(*req.Email, true)
		if err != nil {
			return nil, err
		}

		if email != user.Email {

			existing, err := s.userRepo.FindByEmail(ctx, email)
			if err != nil {
				return nil, err
			}

			if existing != nil && existing.ID != user.ID {
				return nil, errors.New("email already in use with another user")
			}
		}

		user.Email = email
	}

	// Phone
	if req.Phone != nil {

		phone, err := validation.ValidatePhone(*req.Phone)
		if err != nil {
			return nil, err
		}

		if phone != user.Phone {

			existing, err := s.userRepo.FindByPhone(ctx, phone)
			if err != nil {
				return nil, err
			}

			if existing != nil && existing.ID != user.ID {
				return nil, errors.New("phone number already in use with another user")
			}
		}

		user.Phone = phone
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
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

	var err error

	req.FirstName, err = validation.ValidateName(req.FirstName)
	if err != nil {
		return nil, err
	}

	if req.LastName != "" {
		req.LastName, err = validation.ValidateName(req.LastName)
		if err != nil {
			return nil, err
		}
	}

	req.Phone, err = validation.ValidatePhone(req.Phone)
	if err != nil {
		return nil, err
	}

	req.Email, err = validation.ValidateEmail(req.Email, true)
	if err != nil {
		return nil, err
	}

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

	_ = s.auditService.Log(
		ctx,
		models.EntityCustomer,
		user.ID,
		models.AuditCustomerCreated,
		creator.ID,
		bson.M{
			"phone": user.Phone,
			"email": user.Email,
		},
	)
	resp := dto.ToCustomerResponse(user)

	return &resp, nil
}

func (s *UserService) SetCustomerPassword(
	ctx context.Context,
	customerID string,
	updatedBy string,
	req dto.SetCustomerPasswordRequest,
) error {

	id, err := bson.ObjectIDFromHex(customerID)
	if err != nil {
		return errors.New("invalid customer id")
	}

	customer, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if customer == nil {
		return errors.New("customer not found")
	}

	if customer.Role != models.RoleCustomer {
		return errors.New("invalid customer")
	}

	updatedByID, err := bson.ObjectIDFromHex(updatedBy)
	if err != nil {
		return errors.New("invalid user")
	}

	if strings.TrimSpace(req.Password) == "" {
		return errors.New("password is required")
	}

	if err := validation.ValidatePassword(req.Password); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	customer.PasswordHash = string(hash)
	customer.UpdatedAt = time.Now()
	customer.UpdatedByID = &updatedByID

	if err := s.userRepo.UpdatePassword(ctx, customer); err != nil {
		return err
	}
	_ = s.auditService.Log(
		ctx,
		models.EntityUser,
		customer.ID,
		models.AuditPasswordReset,
		updatedByID,
		bson.M{
			"target_user_email": customer.Email,
			"target_role":       customer.Role,
			"reset_by_role":     updatedByID,
		},
	)

	return nil
}

func (s *UserService) UpdateCustomer(
	ctx context.Context,
	customerID string,
	updatedBy string,
	req dto.UpdateCustomerRequest,
) (*dto.CustomerResponse, error) {

	customerObjectID, err := bson.ObjectIDFromHex(customerID)
	if err != nil {
		return nil, errors.New("invalid customer id")
	}

	updatedByID, err := bson.ObjectIDFromHex(updatedBy)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	customer, err := s.userRepo.FindByID(ctx, customerObjectID)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, errors.New("customer not found")
	}

	if customer.Role != models.RoleCustomer {
		return nil, errors.New("invalid customer")
	}

	updater, err := s.userRepo.FindByID(ctx, updatedByID)
	if err != nil {
		return nil, err
	}

	if updater == nil {
		return nil, errors.New("user not found")
	}

	// First Name
	if req.FirstName != nil {

		firstName, err := validation.ValidateName(*req.FirstName)
		if err != nil {
			return nil, err
		}

		customer.FirstName = firstName
	}

	// Last Name
	if req.LastName != nil {

		lastName, err := validation.ValidateName(*req.LastName)
		if err != nil {
			return nil, err
		}

		customer.LastName = lastName
	}

	// Phone
	if req.Phone != nil {

		phone, err := validation.ValidatePhone(*req.Phone)
		if err != nil {
			return nil, err
		}

		existingPhone, err := s.userRepo.FindByPhone(ctx, phone)
		if err != nil {
			return nil, err
		}

		if existingPhone != nil && existingPhone.ID != customer.ID {
			return nil, errors.New("phone number already exists")
		}

		customer.Phone = phone
	}

	// Email
	if req.Email != nil {

		email, err := validation.ValidateEmail(*req.Email, true)
		if err != nil {
			return nil, err
		}

		if email != "" {
			existingEmail, err := s.userRepo.FindByEmail(ctx, email)
			if err != nil {
				return nil, err
			}

			if existingEmail != nil && existingEmail.ID != customer.ID {
				return nil, errors.New("email already exists")
			}
		}

		customer.Email = email
	}

	customer.UpdatedAt = time.Now()
	customer.UpdatedByID = &updatedByID

	if err := s.userRepo.Update(ctx, customer); err != nil {
		return nil, err
	}

	_ = s.auditService.Log(
		ctx,
		models.EntityCustomer,
		customer.ID,
		models.AuditCustomerUpdated,
		updater.ID,
		bson.M{
			"phone": customer.Phone,
			"email": customer.Email,
		},
	)
	resp := dto.ToCustomerResponse(customer)

	return &resp, nil
}
