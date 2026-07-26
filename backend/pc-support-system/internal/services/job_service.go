package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobService struct {
	jobRepo    repository.JobRepository
	userRepo   repository.UserRepository
	deviceRepo repository.DeviceRepository
}

func NewJobService(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	deviceRepo repository.DeviceRepository,
) *JobService {

	return &JobService{
		jobRepo:    jobRepo,
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
	}
}

func (s *JobService) CreateJob(
	ctx context.Context,
	createdBy string,
	req dto.CreateJobRequest,
) (*dto.JobResponse, error) {

	// Customer ID
	customerID, err := bson.ObjectIDFromHex(req.CustomerID)
	if err != nil {
		return nil, errors.New("invalid customer id")
	}

	// Device ID
	deviceID, err := bson.ObjectIDFromHex(req.DeviceID)
	if err != nil {
		return nil, errors.New("invalid device id")
	}

	// Receptionist ID
	createdByID, err := bson.ObjectIDFromHex(createdBy)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	// Customer
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

	// Device
	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	if device == nil {
		return nil, errors.New("device not found")
	}

	// Device belongs to customer
	if device.CustomerID != customerID {
		return nil, errors.New("device does not belong to customer")
	}

	// Receptionist
	creator, err := s.userRepo.FindByID(ctx, createdByID)
	if err != nil {
		return nil, err
	}

	if creator == nil {
		return nil, errors.New("creator not found")
	}

	switch creator.Role {
	case models.RoleReceptionist,
		models.RoleAdmin,
		models.RoleSuperAdmin:
		// allowed
	default:
		return nil, errors.New("user cannot create jobs")
	}

	if creator.State != models.UserActive {
		return nil, errors.New("creator account is inactive")
	}

	// Generate Job Number
	jobNumber, err := s.jobRepo.GenerateJobNumber(ctx)
	if err != nil {
		return nil, err
	}

	req.ProblemDescription = strings.TrimSpace(req.ProblemDescription)

	if req.ProblemDescription == "" {
		return nil, errors.New("problem description is required")
	}

	now := time.Now()

	job := &models.Job{
		JobNumber:          jobNumber,
		CustomerID:         customerID,
		DeviceID:           deviceID,
		ProblemDescription: req.ProblemDescription,

		Status: models.JobCreated,

		CreatedByID: createdByID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = s.jobRepo.Create(ctx, job)
	if err != nil {
		return nil, err
	}

	resp := dto.JobResponse{
		ID:        job.ID.Hex(),
		JobNumber: job.JobNumber,
		Status:    job.Status,

		ProblemDescription: job.ProblemDescription,
		CreatedAt:          job.CreatedAt,

		Customer: dto.CustomerSummary{
			ID:        customer.ID.Hex(),
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Phone:     customer.Phone,
		},

		Device: dto.DeviceSummary{
			ID:    device.ID.Hex(),
			Type:  device.Type,
			Brand: device.Brand,
			Model: device.Model,
		},

		CreatedBy: dto.UserSummary{
			ID:        creator.ID.Hex(),
			FirstName: creator.FirstName,
			LastName:  creator.LastName,
			Role:      creator.Role,
		},
	}

	return &resp, nil
}
