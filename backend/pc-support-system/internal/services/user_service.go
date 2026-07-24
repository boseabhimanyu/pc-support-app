package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
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
