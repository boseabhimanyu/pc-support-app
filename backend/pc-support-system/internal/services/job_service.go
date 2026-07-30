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

	return s.buildJobSummaryResponse(ctx, job)
}

func (s *JobService) GetOpenJobs(
	ctx context.Context,
) (*dto.OpenJobsResponse, error) {

	jobs, err := s.jobRepo.FindByStatuses(
		ctx,
		[]models.JobStatus{
			models.JobCreated,
		},
	)

	if err != nil {
		return nil, err
	}

	resp := make([]dto.JobResponse, 0, len(jobs))

	for _, job := range jobs {

		jobResp, err := s.buildJobSummaryResponse(ctx, job)
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

	jobs, err := s.jobRepo.FindByStatuses(
		ctx,
		[]models.JobStatus{
			models.JobAssigned,
		},
	)

	if err != nil {
		return nil, err
	}

	resp := make([]dto.JobResponse, 0, len(jobs))

	for _, job := range jobs {

		jobResp, err := s.buildJobSummaryResponse(ctx, job)
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

func (s *JobService) GetInProgressJobs(
	ctx context.Context,
) (*dto.JobQueueResponse, error) {

	jobs, err := s.jobRepo.FindByStatuses(
		ctx,
		[]models.JobStatus{
			models.JobInProgress,
		},
	)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.JobResponse, 0, len(jobs))

	for _, job := range jobs {

		jobResp, err := s.buildJobSummaryResponse(ctx, job)
		if err != nil {
			return nil, err
		}

		resp = append(resp, *jobResp)
	}

	return &dto.JobQueueResponse{
		JobsCount: len(resp),
		Jobs:      resp,
	}, nil
}

func (s *JobService) GetWaitingCustomerJobs(
	ctx context.Context,
) (*dto.JobQueueResponse, error) {

	jobs, err := s.jobRepo.FindByStatuses(
		ctx,
		[]models.JobStatus{
			models.JobWaitingCustomer,
		},
	)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.JobResponse, 0, len(jobs))

	for _, job := range jobs {

		jobResp, err := s.buildJobSummaryResponse(ctx, job)
		if err != nil {
			return nil, err
		}

		resp = append(resp, *jobResp)
	}

	return &dto.JobQueueResponse{
		JobsCount: len(resp),
		Jobs:      resp,
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

	// Every assignment resets the workflow
	job.Status = models.JobAssigned
	job.UpdatedAt = time.Now()

	if err := s.jobRepo.AssignJob(ctx, job); err != nil {
		return nil, err
	}

	return s.buildJobSummaryResponse(ctx, job)
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

		jobResp, err := s.buildJobSummaryResponse(ctx, job)
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

func (s *JobService) buildJobSummaryResponse(
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

		jobResp, err := s.buildJobSummaryResponse(ctx, job)
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

func (s *JobService) ChangeJobStatus(
	ctx context.Context,
	jobID string,
	userID string,
	req dto.UpdateJobStatusRequest,
) (*dto.JobResponse, error) {

	jobObjectID, err := bson.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, errors.New("invalid job")
	}

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	job, err := s.jobRepo.FindByID(ctx, jobObjectID)
	if err != nil {
		return nil, err
	}

	if job == nil {
		return nil, errors.New("job not found")
	}

	staff, err := s.userRepo.FindByID(ctx, userObjectID)
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
		return nil, errors.New("user cannot update job status")
	}

	if staff.State != models.UserActive {
		return nil, errors.New("user account is inactive")
	}

	if job.AssignedToID == nil {
		return nil, errors.New("job is not assigned")
	}

	if *job.AssignedToID != staff.ID {
		return nil, errors.New("job is not assigned to you")
	}

	if job.Status == models.JobClosed {
		return nil, errors.New("closed jobs cannot be updated")
	}

	if !req.Status.IsValid() {
		return nil, errors.New("invalid job status")
	}

	if req.Status == models.JobClosed {
		return nil, errors.New("jobs must be closed using the close endpoint")
	}

	switch job.Status {

	case models.JobAssigned:
		if req.Status != models.JobInProgress {
			return nil, errors.New("assigned jobs can only move to in progress")
		}

	case models.JobInProgress:
		if req.Status != models.JobWaitingCustomer {
			return nil, errors.New("invalid status transition")
		}

	case models.JobWaitingCustomer:
		if req.Status != models.JobResumed {
			return nil, errors.New("invalid status transition")
		}

	case models.JobResumed:
		if req.Status != models.JobWaitingCustomer {
			return nil, errors.New("invalid status transition")
		}

	default:
		return nil, errors.New("invalid current job status")
	}

	err = s.jobRepo.UpdateStatus(
		ctx,
		job.ID,
		req.Status,
	)
	if err != nil {
		return nil, err
	}

	job.Status = req.Status
	job.UpdatedAt = time.Now()

	return s.buildJobSummaryResponse(ctx, job)
}

func (s *JobService) GetResumedJobs(
	ctx context.Context,
) (*dto.JobQueueResponse, error) {

	jobs, err := s.jobRepo.FindByStatuses(
		ctx,
		[]models.JobStatus{
			models.JobResumed,
		},
	)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.JobResponse, 0, len(jobs))

	for _, job := range jobs {

		jobResp, err := s.buildJobSummaryResponse(ctx, job)
		if err != nil {
			return nil, err
		}

		resp = append(resp, *jobResp)
	}

	return &dto.JobQueueResponse{
		JobsCount: len(resp),
		Jobs:      resp,
	}, nil
}

func (s *JobService) AddJobNote(
	ctx context.Context,
	jobID string,
	userID string,
	req dto.AddJobNoteRequest,
) error {

	req.Note = strings.TrimSpace(req.Note)

	if req.Note == "" {
		return errors.New("note is required")
	}

	if len(req.Note) > 2000 {
		return errors.New("note cannot exceed 2000 characters")
	}

	jobObjectID, err := bson.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New("invalid job id")
	}

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user")
	}

	job, err := s.jobRepo.FindByID(ctx, jobObjectID)
	if err != nil {
		return err
	}

	if job == nil {
		return errors.New("job not found")
	}

	user, err := s.userRepo.FindByID(ctx, userObjectID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if user.State != models.UserActive {
		return errors.New("user account is inactive")
	}

	// Closed jobs cannot receive notes
	if job.Status == models.JobClosed {
		return errors.New("cannot add notes to a closed job")
	}

	switch user.Role {

	case models.RoleReceptionist:
		// Receptionists can add customer communication notes until the job is closed.
	case models.RoleTechnician, models.RoleHeadTechnician:

		if job.AssignedToID == nil {
			return errors.New("job is not assigned")
		}

		if *job.AssignedToID != user.ID {
			return errors.New("job is not assigned to you")
		}

	default:
		return errors.New("user cannot add job notes")
	}

	note := models.JobNote{
		ID:         bson.NewObjectID(),
		AuthorID:   user.ID,
		AuthorRole: user.Role,
		Note:       req.Note,
		CreatedAt:  time.Now(),
	}

	return s.jobRepo.AddNote(
		ctx,
		job.ID,
		note,
	)
}

func (s *JobService) GetJobNotes(
	ctx context.Context,
	jobID string,
	userID string,
) (*dto.JobNotesResponse, error) {

	jobObjectID, err := bson.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, errors.New("invalid job id")
	}

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	job, err := s.jobRepo.FindByID(ctx, jobObjectID)
	if err != nil {
		return nil, err
	}

	if job == nil {
		return nil, errors.New("job not found")
	}

	user, err := s.userRepo.FindByID(ctx, userObjectID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	switch user.Role {

	case models.RoleReceptionist,
		models.RoleAdmin,
		models.RoleSuperAdmin:

		// Can view all jobs

	case models.RoleHeadTechnician:

		// Can view all jobs

	case models.RoleTechnician:

		if job.AssignedToID == nil || *job.AssignedToID != user.ID {
			return nil, errors.New("access denied")
		}

	case models.RoleCustomer:

		if job.CustomerID != user.ID {
			return nil, errors.New("access denied")
		}

	default:
		return nil, errors.New("access denied")
	}

	resp := make([]dto.JobNoteResponse, 0, len(job.Notes))

	for _, note := range job.Notes {

		author, _ := s.userRepo.FindByID(ctx, note.AuthorID)

		if author == nil {
			continue
		}

		resp = append(resp, dto.ToJobNoteResponse(
			note,
			author,
		))
	}

	return &dto.JobNotesResponse{
		NotesCount: len(resp),
		Notes:      resp,
	}, nil
}

func (s *JobService) GetJobByID(
	ctx context.Context,
	jobID string,
	userID string,
) (*dto.JobResponse, error) {

	jobObjectID, err := bson.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, errors.New("invalid job id")
	}

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	job, err := s.jobRepo.FindByID(ctx, jobObjectID)
	if err != nil {
		return nil, err
	}

	if job == nil {
		return nil, errors.New("job not found")
	}

	user, err := s.userRepo.FindByID(ctx, userObjectID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.State != models.UserActive {
		return nil, errors.New("user account is inactive")
	}

	switch user.Role {

	case models.RoleReceptionist,
		models.RoleTechnician,
		models.RoleHeadTechnician,
		models.RoleAdmin,
		models.RoleSuperAdmin:
		// Staff can view any job.

	case models.RoleCustomer:

		if job.CustomerID != user.ID {
			return nil, errors.New("access denied")
		}

	default:
		return nil, errors.New("access denied")
	}

	resp, err := s.buildJobDetailsResponse(ctx, job)
	if err != nil {
		return nil, err
	}

	if user.Role == models.RoleCustomer {
		resp.CloseReason = ""
		resp.InternalClosureNotes = ""
	}

	return resp, nil
}

func (s *JobService) buildJobDetailsResponse(
	ctx context.Context,
	job *models.Job,
) (*dto.JobResponse, error) {

	resp, err := s.buildJobSummaryResponse(ctx, job)
	if err != nil {
		return nil, err
	}

	notes := make([]dto.JobNoteResponse, 0, len(job.Notes))

	for _, note := range job.Notes {

		author, err := s.userRepo.FindByID(ctx, note.AuthorID)
		if err != nil || author == nil {
			continue
		}

		notes = append(notes, dto.ToJobNoteResponse(
			note,
			author,
		))
	}

	resp.Notes = notes

	return resp, nil
}

func (s *JobService) GetCustomerJobsByCustomerID(
	ctx context.Context,
	customerID string,
	userID string,
) (*dto.CustomerJobsResponse, error) {

	customerObjectID, err := bson.ObjectIDFromHex(customerID)
	if err != nil {
		return nil, errors.New("invalid customer id")
	}

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	user, err := s.userRepo.FindByID(ctx, userObjectID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.State != models.UserActive {
		return nil, errors.New("user account is inactive")
	}

	switch user.Role {
	case models.RoleReceptionist,
		models.RoleTechnician,
		models.RoleHeadTechnician,
		models.RoleAdmin,
		models.RoleSuperAdmin:
		// allowed

	default:
		return nil, errors.New("user cannot view customer jobs")
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

	jobs, err := s.jobRepo.FindCustomerJobs(
		ctx,
		customerObjectID,
	)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.JobResponse, 0, len(jobs))

	for _, job := range jobs {

		jobResp, err := s.buildJobSummaryResponse(
			ctx,
			job,
		)
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

func (s *JobService) CloseJob(
	ctx context.Context,
	jobID string,
	userID string,
	req dto.CloseJobRequest,
) (*dto.JobResponse, error) {

	jobObjectID, err := bson.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, errors.New("invalid job id")
	}

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	job, err := s.jobRepo.FindByID(ctx, jobObjectID)
	if err != nil {
		return nil, err
	}

	if job == nil {
		return nil, errors.New("job not found")
	}

	user, err := s.userRepo.FindByID(ctx, userObjectID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.State != models.UserActive {
		return nil, errors.New("user account is inactive")
	}

	switch user.Role {

	case models.RoleTechnician,
		models.RoleHeadTechnician:

		// allowed

	default:
		return nil, errors.New("user cannot close jobs")
	}

	if job.AssignedToID == nil {
		return nil, errors.New("job is not assigned")
	}

	if *job.AssignedToID != user.ID {
		return nil, errors.New("job is not assigned to you")
	}

	if job.Status == models.JobClosed {
		return nil, errors.New("job is already closed")
	}

	if job.Status != models.JobResumed {
		return nil, errors.New("only resumed jobs can be closed")
	}

	if !req.Reason.IsValid() {
		return nil, errors.New("invalid closure reason")
	}

	req.ClosureNotes = strings.TrimSpace(req.ClosureNotes)
	req.InternalClosureNotes = strings.TrimSpace(req.InternalClosureNotes)

	now := time.Now()

	job.Status = models.JobClosed
	job.CloseReason = req.Reason
	job.ClosureNotes = req.ClosureNotes
	job.InternalClosureNotes = req.InternalClosureNotes

	job.ClosedByID = user.ID
	job.ClosedAt = &now
	job.UpdatedAt = now

	if err := s.jobRepo.CloseJob(
		ctx,
		job,
	); err != nil {
		return nil, err
	}

	return s.buildJobDetailsResponse(
		ctx,
		job,
	)
}
