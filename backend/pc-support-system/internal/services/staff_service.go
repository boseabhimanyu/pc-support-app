package services

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) CreateStaff(
	ctx context.Context,
	createdBy string,
	req dto.CreateStaffRequest,
) (*dto.UserResponse, error) {

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

	// Email must be unique

	if req.Email != "" {

		existingEmail, err := s.userRepo.FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}

		if existingEmail != nil {
			return nil, errors.New("email already exists")
		}
	}

	switch req.Role {
	case models.RoleReceptionist,
		models.RoleTechnician,
		models.RoleHeadTechnician,
		models.RoleAdmin:
		// valid
	default:
		return nil, errors.New("invalid staff role")
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

	if creator.State != models.UserActive {
		return nil, errors.New("creator account is inactive")
	}

	if creator.Role == models.RoleAdmin &&
		req.Role == models.RoleAdmin {
		return nil, errors.New("admin cannot create another admin")
	}

	switch creator.Role {
	case models.RoleAdmin,
		models.RoleSuperAdmin:
		// allowed
	default:
		return nil, errors.New("user cannot create staff")
	}
	now := time.Now()

	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Email:     req.Email,

		Role: req.Role,

		State: models.UserActive,

		PasswordHash: "",

		CreatedByID: &createdByID,

		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	resp := dto.ToUserResponse(user)

	return &resp, nil

}

func (s *UserService) SetStaffPassword(
	ctx context.Context,
	staffID string,
	updatedBy string,
	req dto.SetStaffPasswordRequest,
) error {

	// Staff ID
	id, err := bson.ObjectIDFromHex(staffID)
	if err != nil {
		return errors.New("invalid staff id")
	}

	// Updater ID
	updatedByID, err := bson.ObjectIDFromHex(updatedBy)
	if err != nil {
		return errors.New("invalid user")
	}

	// Staff
	staff, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if staff == nil {
		return errors.New("staff not found")
	}

	// Must not be a customer
	if staff.Role == models.RoleCustomer {
		return errors.New("invalid staff")
	}

	// Updater
	updater, err := s.userRepo.FindByID(ctx, updatedByID)
	if err != nil {
		return err
	}

	if updater == nil {
		return errors.New("user not found")
	}

	if updater.State != models.UserActive {
		return errors.New("user account is inactive")
	}

	// Only Admin / Super Admin can set staff passwords
	switch updater.Role {
	case models.RoleAdmin,
		models.RoleSuperAdmin:
		// allowed
	default:
		return errors.New("user cannot set staff password")
	}

	// Admin cannot reset another Admin's password
	if updater.Role == models.RoleAdmin &&
		staff.Role == models.RoleAdmin {
		return errors.New("admin cannot reset another admin's password")
	}

	// Validate password
	if strings.TrimSpace(req.Password) == "" {
		return errors.New("password is required")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	staff.PasswordHash = string(hash)
	staff.UpdatedAt = time.Now()
	staff.UpdatedByID = &updatedByID

	return s.userRepo.UpdatePassword(ctx, staff)
}

func (s *UserService) SearchStaff(
	ctx context.Context,
	query string,
) ([]dto.UserResponse, error) {

	query = strings.TrimSpace(query)

	staff, err := s.userRepo.SearchStaff(ctx, query)
	if err != nil {
		return nil, err
	}

	response := make([]dto.UserResponse, 0, len(staff))

	for _, user := range staff {

		response = append(
			response,
			dto.ToUserResponse(&user),
		)
	}

	return response, nil
}
