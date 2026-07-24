package services

import (
	"context"
	"errors"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewDeviceService(
	deviceRepo repository.DeviceRepository,
	userRepo repository.UserRepository,
) *DeviceService {

	return &DeviceService{
		deviceRepo: deviceRepo,
		userRepo:   userRepo,
	}
}

type DeviceService struct {
	deviceRepo repository.DeviceRepository
	userRepo   repository.UserRepository
}

func (s *DeviceService) AddDevice(
	ctx context.Context,
	req dto.AddDeviceRequest,
) (*dto.DeviceResponse, error) {

	// Convert CustomerID to ObjectID
	customerID, err := bson.ObjectIDFromHex(req.CustomerID)
	if err != nil {
		return nil, errors.New("invalid customer id")
	}

	// Check customer exists
	user, err := s.userRepo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("customer not found")
	}

	// Ensure the user is a customer
	if user.Role != models.RoleCustomer {
		return nil, errors.New("device can only be assigned to a customer")
	}

	// Create device
	device := &models.Device{
		CustomerID: customerID,
		Type:       req.Type,

		Brand:        req.Brand,
		Model:        req.Model,
		SerialNumber: req.SerialNumber,
		Notes:        req.Notes,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save
	err = s.deviceRepo.Create(ctx, device)
	if err != nil {
		return nil, err
	}

	resp := dto.ToDeviceResponse(device)

	return &resp, nil
}

func (s *DeviceService) GetCustomerDevices(
	ctx context.Context,
	customerID string,
) ([]dto.DeviceResponse, error) {

	id, err := bson.ObjectIDFromHex(customerID)
	if err != nil {
		return nil, errors.New("invalid customer id")
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("customer not found")
	}

	if user.Role != models.RoleCustomer {
		return nil, errors.New("invalid customer")
	}

	devices, err := s.deviceRepo.FindByCustomerID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToDeviceResponses(devices), nil
}

// UpdateDevice()

// DeactivateDevice()
