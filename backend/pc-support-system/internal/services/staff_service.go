package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/validation"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) CreateStaff(
	ctx context.Context,
	createdBy string,
	req dto.CreateStaffRequest,
) (*dto.UserResponse, error) {

	var err error

	// Validate & normalize input
	req.FirstName, err = validation.ValidateName(req.FirstName)
	if err != nil {
		return nil, err
	}

	req.LastName, err = validation.ValidateName(req.LastName)
	if err != nil {
		return nil, err
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

	// Validate role
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

	switch creator.Role {
	case models.RoleAdmin,
		models.RoleSuperAdmin:
		// allowed
	default:
		return nil, errors.New("user cannot create staff")
	}

	if creator.Role == models.RoleAdmin &&
		req.Role == models.RoleAdmin {
		return nil, errors.New("admin cannot create another admin")
	}

	now := time.Now()

	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Email:     req.Email,

		Role:  req.Role,
		State: models.UserActive,

		PasswordHash: "",

		CreatedByID: &createdByID,

		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	_ = s.auditService.Log(
		ctx,
		models.EntityUser,
		user.ID,
		models.AuditUserCreated,
		creator.ID,
		bson.M{
			"role":  user.Role,
			"state": user.State,
			"phone": user.Phone,
			"email": user.Email,
		},
	)
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

	if err := s.userRepo.UpdatePassword(ctx, staff); err != nil {
		return err
	}

	_ = s.auditService.Log(
		ctx,
		models.EntityUser,
		staff.ID,
		models.AuditPasswordChanged,
		updatedByID,
		nil,
	)

	return nil
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

func (s *UserService) UpdateStaff(
	ctx context.Context,
	staffID string,
	updatedBy string,
	req dto.UpdateStaffRequest,
) (*dto.UserResponse, error) {

	// Staff ID
	staffObjectID, err := bson.ObjectIDFromHex(staffID)
	if err != nil {
		return nil, errors.New("invalid staff id")
	}

	// Updater ID
	updatedByID, err := bson.ObjectIDFromHex(updatedBy)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	// Staff
	staff, err := s.userRepo.FindByID(ctx, staffObjectID)
	if err != nil {
		return nil, err
	}

	if staff == nil {
		return nil, errors.New("staff not found")
	}

	if staff.Role == models.RoleCustomer {
		return nil, errors.New("invalid staff")
	}

	// Updater
	updater, err := s.userRepo.FindByID(ctx, updatedByID)
	if err != nil {
		return nil, err
	}

	if updater == nil {
		return nil, errors.New("user not found")
	}

	if updater.State != models.UserActive {
		return nil, errors.New("user account is inactive")
	}

	switch updater.Role {
	case models.RoleAdmin,
		models.RoleSuperAdmin:
		// allowed
	default:
		return nil, errors.New("user cannot update staff")
	}

	// First Name
	if req.FirstName != nil {

		firstName, err := validation.ValidateName(*req.FirstName)
		if err != nil {
			return nil, err
		}

		staff.FirstName = firstName
	}

	// Last Name
	if req.LastName != nil {

		lastName, err := validation.ValidateName(*req.LastName)
		if err != nil {
			return nil, err
		}

		staff.LastName = lastName
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

		if existingPhone != nil && existingPhone.ID != staff.ID {
			return nil, errors.New("phone number already exists")
		}

		staff.Phone = phone
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

			if existingEmail != nil && existingEmail.ID != staff.ID {
				return nil, errors.New("email already exists")
			}
		}

		staff.Email = email
	}

	staff.UpdatedAt = time.Now()
	staff.UpdatedByID = &updatedByID

	if err := s.userRepo.Update(ctx, staff); err != nil {
		return nil, err
	}

	_ = s.auditService.Log(
		ctx,
		models.EntityUser,
		staff.ID,
		models.AuditUserUpdated,
		updater.ID,
		bson.M{
			"role":  staff.Role,
			"state": staff.State,
			"phone": staff.Phone,
			"email": staff.Email,
		},
	)
	resp := dto.ToUserResponse(staff)

	return &resp, nil
}
