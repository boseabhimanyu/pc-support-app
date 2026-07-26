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

type UserSummary struct {
	ID        string      `json:"id"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Role      models.Role `json:"role"`
}
