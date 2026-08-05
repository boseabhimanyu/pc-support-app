package dto

import (
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
)

type CreateJobRequest struct {
	CustomerID string `json:"customerId" binding:"required"`
	DeviceID   string `json:"deviceId" binding:"required"`

	ProblemDescription string `json:"problemDescription" binding:"required"`
}

type JobResponse struct {
	ID        string `json:"id"`
	JobNumber string `json:"jobNumber"`

	Status models.JobStatus `json:"status"`

	Customer CustomerSummary `json:"customer"`
	Device   DeviceSummary   `json:"device"`

	ProblemDescription string `json:"problemDescription"`

	Notes []JobNoteResponse `json:"notes"`

	CreatedAt time.Time `json:"createdAt"`

	CreatedBy UserSummary `json:"createdBy"`

	AssignedTo *UserSummary `json:"assignedTo,omitempty"`

	CloseReason          models.JobCloseReason `json:"closeReason,omitempty"`
	ClosureNotes         string                `json:"closureNotes,omitempty"`
	ClosedAt             *time.Time            `json:"closedAt,omitempty"`
	InternalClosureNotes string                `json:"internalClosureNotes,omitempty"`
}

type DeviceSummary struct {
	ID string `json:"id"`

	Type models.DeviceType `json:"type"`

	Brand string `json:"brand,omitempty"`

	Model string `json:"model,omitempty"`

	SerialNumber string `json:"serialNumber,omitempty"`
}

type AssignTechnicianRequest struct {
	TechnicianID string `json:"technicianId" binding:"required"`
}

func ToJobResponse(
	job *models.Job,
	customer CustomerSummary,
	device DeviceSummary,
	createdBy UserSummary,
	assignedTo *UserSummary,
) JobResponse {

	return JobResponse{
		ID:                 job.ID.Hex(),
		JobNumber:          job.JobNumber,
		Status:             job.Status,
		Customer:           customer,
		Device:             device,
		ProblemDescription: job.ProblemDescription,
		CreatedAt:          job.CreatedAt,
		CreatedBy:          createdBy,
		AssignedTo:         assignedTo,

		CloseReason:          job.CloseReason,
		ClosureNotes:         job.ClosureNotes,
		InternalClosureNotes: job.InternalClosureNotes,
		ClosedAt:             job.ClosedAt,
	}
}

type OpenJobResponse struct {
	ID         string           `json:"id"`
	JobNumber  string           `json:"jobNumber"`
	CustomerID string           `json:"customerId"`
	DeviceID   string           `json:"deviceId"`
	Status     models.JobStatus `json:"status"`
	CreatedAt  time.Time        `json:"createdAt"`
}

type OpenJobsResponse struct {
	OpenJobsCount int           `json:"openJobsCount"`
	Jobs          []JobResponse `json:"jobs"`
}

type AssignJobRequest struct {
	StaffID string `json:"staffId" binding:"required"`
	//Note    string `json:"note"`
}

type AssignedJobsResponse struct {
	AssisgnedJobsCount int           `json:"assignedJobsCount"`
	Jobs               []JobResponse `json:"jobs"`
}

type MyJobsResponse struct {
	MyJobsCount int           `json:"myJobsCount"`
	Jobs        []JobResponse `json:"jobs"`
}

type CustomerJobsResponse struct {
	JobsCount int           `json:"jobsCount"`
	Jobs      []JobResponse `json:"jobs"`
}

type UpdateJobStatusRequest struct {
	Status models.JobStatus `json:"status" binding:"required"`
}

type JobQueueResponse struct {
	JobsCount int           `json:"jobsCount"`
	Jobs      []JobResponse `json:"jobs"`
}

type AddJobNoteRequest struct {
	Note string `json:"note" binding:"required"`
}

type JobNotesResponse struct {
	NotesCount int               `json:"notesCount"`
	Notes      []JobNoteResponse `json:"notes"`
}

type JobNoteResponse struct {
	ID        string      `json:"id"`
	Author    UserSummary `json:"author"`
	Note      string      `json:"note"`
	CreatedAt time.Time   `json:"createdAt"`
}

func ToJobNoteResponse(
	note models.JobNote,
	author *models.User,
) JobNoteResponse {

	return JobNoteResponse{
		ID: note.ID.Hex(),
		Author: UserSummary{
			ID:        author.ID.Hex(),
			FirstName: author.FirstName,
			LastName:  author.LastName,
			Role:      note.AuthorRole, // historical role
		},
		Note:      note.Note,
		CreatedAt: note.CreatedAt,
	}
}

type CloseJobRequest struct {
	Reason               models.JobCloseReason `json:"reason" binding:"required"`
	ClosureNotes         string                `json:"closureNotes"`
	InternalClosureNotes string                `json:"internalClosureNotes"`
}
