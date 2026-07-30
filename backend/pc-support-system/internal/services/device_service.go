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
	userID string,
	req dto.AddDeviceRequest,
) (*dto.DeviceResponse, error) {

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	staff, err := s.userRepo.FindByID(ctx, userObjectID)
	if err != nil {
		return nil, err
	}

	if staff == nil {
		return nil, errors.New("user not found")
	}

	if staff.State != models.UserActive {
		return nil, errors.New("user account is inactive")
	}

	switch staff.Role {
	case models.RoleReceptionist,
		models.RoleHeadTechnician,
		models.RoleAdmin,
		models.RoleSuperAdmin:
		// allowed

	default:
		return nil, errors.New("user cannot add devices")
	}

	// Convert CustomerID to ObjectID
	customerID, err := bson.ObjectIDFromHex(req.CustomerID)
	if err != nil {
		return nil, errors.New("invalid customer id")
	}

	if !req.Type.IsValid() {
		return nil, errors.New("invalid device type")
	}

	if !req.Condition.IsValid() {
		return nil, errors.New("invalid device condition")
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
		Condition:  req.Condition,

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

	resp.Customer = dto.ToCustomerSummary(user)

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

	responses := dto.ToDeviceResponses(devices)

	for i := range responses {
		responses[i].Customer = dto.CustomerSummary{
			ID:        user.ID.Hex(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
		}
	}

	return responses, nil
}

func (s *DeviceService) GetDevice(
	ctx context.Context,
	id string,
) (*dto.DeviceResponse, error) {

	deviceID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid device id")
	}

	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	if device == nil {
		return nil, errors.New("device not found")
	}

	// Get customer details
	user, err := s.userRepo.FindByID(ctx, device.CustomerID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("customer not found")
	}

	resp := dto.ToDeviceResponse(device)

	resp.Customer = dto.CustomerSummary{
		ID:        user.ID.Hex(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
	}

	return &resp, nil
}

func (s *DeviceService) GetMyDevices(
	ctx context.Context,
	userID string,
) (*dto.CustomerDevicesResponse, error) {

	customerID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	customer, err := s.userRepo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, errors.New("customer not found")
	}

	if customer.Role != models.RoleCustomer {
		return nil, errors.New("invalid customer")
	}

	if customer.State != models.UserActive {
		return nil, errors.New("customer account is inactive")
	}

	devices, err := s.deviceRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.CustomerDeviceResponse, 0, len(devices))

	for i := range devices {
		resp = append(resp, dto.ToCustomerDeviceResponse(&devices[i]))
	}

	return &dto.CustomerDevicesResponse{
		DevicesCount: len(resp),
		Devices:      resp,
	}, nil
}

// UpdateDevice()

// DeactivateDevice()
