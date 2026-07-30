package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobStatus string

const (
	JobCreated         JobStatus = "created"
	JobAssigned        JobStatus = "assigned"
	JobInProgress      JobStatus = "in_progress"
	JobWaitingCustomer JobStatus = "waiting_customer"
	JobResumed         JobStatus = "resumed"
	JobClosed          JobStatus = "closed"
)

func (s JobStatus) IsValid() bool {
	switch s {
	case JobCreated,
		JobAssigned,
		JobInProgress,
		JobWaitingCustomer,
		JobResumed,
		JobClosed:
		return true
	default:
		return false
	}
}

type Job struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`

	// Human-readable Job ID
	JobNumber string `bson:"job_number" json:"jobNumber"`

	// Relationships
	CustomerID bson.ObjectID `bson:"customer_id" json:"customerId"`
	DeviceID   bson.ObjectID `bson:"device_id" json:"deviceId"`

	// Job Status
	Status JobStatus `bson:"status" json:"status"`

	// Original customer complaint (immutable)
	ProblemDescription string `bson:"problem_description" json:"problemDescription"`

	// Closure
	CloseReason JobCloseReason `bson:"close_reason,omitempty" json:"closeReason,omitempty"`

	// Visible to customer
	ClosureNotes string `bson:"closure_notes,omitempty" json:"closureNotes,omitempty"`

	// Internal staff only
	InternalClosureNotes string `bson:"internal_closure_notes,omitempty" json:"-"`

	// Staff
	CreatedByID bson.ObjectID `bson:"created_by_id" json:"createdById"`

	AssignedToID *bson.ObjectID `bson:"assigned_to_id,omitempty" json:"assignedToId,omitempty"`
	AssignedByID *bson.ObjectID `bson:"assigned_by_id,omitempty" json:"assignedById,omitempty"`

	ClosedByID bson.ObjectID `bson:"closed_by_id,omitempty" json:"closedById,omitempty"`

	// Timeline
	Notes []JobNote `bson:"notes,omitempty" json:"notes,omitempty"`

	CreatedAt time.Time  `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updatedAt"`
	ClosedAt  *time.Time `bson:"closed_at,omitempty" json:"closedAt,omitempty"`
}

type JobNote struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`

	AuthorID bson.ObjectID `bson:"author_id" json:"authorId"`

	// Role when the note was written
	AuthorRole Role `bson:"author_role" json:"authorRole"`

	Note string `bson:"note" json:"note"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
}

type JobCloseReason string

const (
	JobCompleted          JobCloseReason = "completed"
	JobNotRepairable      JobCloseReason = "not_repairable"
	JobCustomerCancelled  JobCloseReason = "customer_cancelled"
	JobCustomerNoResponse JobCloseReason = "customer_no_response"
	JobDuplicateJob       JobCloseReason = "duplicate_job"
	JobWrongSerial        JobCloseReason = "wrong_serial"
)

func (r JobCloseReason) IsValid() bool {
	switch r {
	case JobCompleted,
		JobNotRepairable,
		JobCustomerCancelled,
		JobCustomerNoResponse,
		JobDuplicateJob,
		JobWrongSerial:
		return true
	default:
		return false
	}
}
