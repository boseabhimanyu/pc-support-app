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
		models.RoleHeadTechnician,
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

	return s.buildJobResponse(ctx, job)
}

func (s *JobService) GetOpenJobs(
	ctx context.Context,
) (*dto.OpenJobsResponse, error) {

	jobs, err := s.jobRepo.FindOpenUnassigned(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.JobResponse, 0, len(jobs))

	for _, job := range jobs {

		jobResp, err := s.buildJobResponse(ctx, job)
		if err != nil {
			return nil, err
		}

		resp = append(resp, *jobResp)
	}

	return &dto.OpenJobsResponse{
		OpenJobsCount: len(resp),
		Jobs:          resp,
	}, nil
}

func (s *JobService) GetAssignedJobs(
	ctx context.Context,
) (*dto.AssignedJobsResponse, error) {

	jobs, err := s.jobRepo.FindOpenAssigned(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.JobResponse, 0, len(jobs))

	for _, job := range jobs {

		jobResp, err := s.buildJobResponse(ctx, job)
		if err != nil {
			return nil, err
		}

		resp = append(resp, *jobResp)
	}

	return &dto.AssignedJobsResponse{
		AssisgnedJobsCount: len(resp),
		Jobs:               resp,
	}, nil
}

func (s *JobService) AssignJob(
	ctx context.Context,
	jobID string,
	assignedBy string,
	req dto.AssignJobRequest,
) (*dto.JobResponse, error) {

	// Job ID
	jobObjectID, err := bson.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, errors.New("invalid job id")
	}

	// Assignee ID
	staffObjectID, err := bson.ObjectIDFromHex(req.StaffID)
	if err != nil {
		return nil, errors.New("invalid staff id")
	}

	// Assigner ID
	assignedByID, err := bson.ObjectIDFromHex(assignedBy)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	// Job
	job, err := s.jobRepo.FindByID(ctx, jobObjectID)
	if err != nil {
		return nil, err
	}

	if job == nil {
		return nil, errors.New("job not found")
	}

	if job.Status == models.JobClosed {
		return nil, errors.New("closed jobs cannot be assigned")
	}

	// Assignee
	staff, err := s.userRepo.FindByID(ctx, staffObjectID)
	if err != nil {
		return nil, err
	}

	if staff == nil {
		return nil, errors.New("staff not found")
	}

	if staff.State != models.UserActive {
		return nil, errors.New("staff account is inactive")
	}

	switch staff.Role {
	case models.RoleTechnician,
		models.RoleHeadTechnician:
		// valid
	default:
		return nil, errors.New("invalid staff role")
	}

	// Assigner
	assigner, err := s.userRepo.FindByID(ctx, assignedByID)
	if err != nil {
		return nil, err
	}

	if assigner == nil {
		return nil, errors.New("user not found")
	}

	if assigner.State != models.UserActive {
		return nil, errors.New("user account is inactive")
	}

	switch assigner.Role {
	case models.RoleHeadTechnician,
		models.RoleAdmin,
		models.RoleSuperAdmin:
		// allowed
	default:
		return nil, errors.New("user cannot assign jobs")
	}

	// Assignment
	job.AssignedToID = &staff.ID
	job.AssignedByID = &assigner.ID

	if job.Status == models.JobCreated {
		job.Status = models.JobAssigned
	}

	job.UpdatedAt = time.Now()

	if err := s.jobRepo.AssignJob(ctx, job); err != nil {
		return nil, err
	}

	return s.buildJobResponse(ctx, job)
}

func (s *JobService) GetMyJobs(
	ctx context.Context,
	userID string,
) (*dto.MyJobsResponse, error) {

	staffID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	staff, err := s.userRepo.FindByID(ctx, staffID)
	if err != nil {
		return nil, err
	}

	if staff == nil {
		return nil, errors.New("user not found")
	}

	switch staff.Role {
	case models.RoleTechnician,
		models.RoleHeadTechnician:
		// allowed
	default:
		return nil, errors.New("user cannot access technician jobs")
	}

	if staff.State != models.UserActive {
		return nil, errors.New("user account is inactive")
	}

	jobs, err := s.jobRepo.FindMyJobs(ctx, staffID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.JobResponse, 0, len(jobs))

	for _, job := range jobs {

		jobResp, err := s.buildJobResponse(ctx, job)
		if err != nil {
			return nil, err
		}

		resp = append(resp, *jobResp)
	}

	return &dto.MyJobsResponse{
		MyJobsCount: len(resp),
		Jobs:        resp,
	}, nil
}

func (s *JobService) buildJobResponse(
	ctx context.Context,
	job *models.Job,
) (*dto.JobResponse, error) {

	customer, err := s.userRepo.FindByID(ctx, job.CustomerID)
	if err != nil {
		return nil, err
	}

	device, err := s.deviceRepo.FindByID(ctx, job.DeviceID)
	if err != nil {
		return nil, err
	}

	createdBy, err := s.userRepo.FindByID(ctx, job.CreatedByID)
	if err != nil {
		return nil, err
	}

	var assignedTo *dto.UserSummary

	if job.AssignedToID != nil {
		user, err := s.userRepo.FindByID(ctx, *job.AssignedToID)
		if err != nil {
			return nil, err
		}

		if user != nil {
			summary := dto.ToUserSummary(user)
			assignedTo = &summary
		}
	}

	resp := dto.ToJobResponse(
		job,
		dto.ToCustomerSummary(customer),
		dto.ToDeviceSummary(device),
		dto.ToUserSummary(createdBy),
		assignedTo,
	)

	return &resp, nil
}

func (s *JobService) GetCustomerJobs(
	ctx context.Context,
	userID string,
) (*dto.CustomerJobsResponse, error) {

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

	jobs, err := s.jobRepo.FindCustomerJobs(ctx, customerID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.JobResponse, 0, len(jobs))

	for _, job := range jobs {

		jobResp, err := s.buildJobResponse(ctx, job)
		if err != nil {
			return nil, err
		}

		resp = append(resp, *jobResp)
	}

	return &dto.CustomerJobsResponse{
		JobsCount: len(resp),
		Jobs:      resp,
	}, nil
}
