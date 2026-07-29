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

	CreatedAt time.Time `json:"createdAt"`

	CreatedBy UserSummary `json:"createdBy"`

	AssignedTo *UserSummary `json:"assignedTo,omitempty"`
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
