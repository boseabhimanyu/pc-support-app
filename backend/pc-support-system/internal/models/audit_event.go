package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuditEntityType string

const (
	EntityJob      AuditEntityType = "job"
	EntityDevice   AuditEntityType = "device"
	EntityCustomer AuditEntityType = "customer"
	EntityUser     AuditEntityType = "user"
)

type AuditEventType string

const (
	AuditJobCreated       AuditEventType = "job_created"
	AuditJobAssigned      AuditEventType = "job_assigned"
	AuditJobStatusChanged AuditEventType = "job_status_changed"
	AuditJobNoteAdded     AuditEventType = "job_note_added"
	AuditJobClosed        AuditEventType = "job_closed"

	AuditDeviceCreated AuditEventType = "device_created"
	AuditDeviceUpdated AuditEventType = "device_updated"

	AuditCustomerCreated AuditEventType = "customer_created"
	AuditCustomerUpdated AuditEventType = "customer_updated"

	AuditUserCreated     AuditEventType = "user_created"
	AuditUserUpdated     AuditEventType = "user_updated"
	AuditPasswordChanged AuditEventType = "password_changed"
	AuditPasswordReset   AuditEventType = "password_reset"
)

type AuditEvent struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	EntityType AuditEntityType `bson:"entity_type"`
	EntityID   bson.ObjectID   `bson:"entity_id"`

	EventType AuditEventType `bson:"event_type"`

	PerformedByID bson.ObjectID `bson:"performed_by_id"`

	Metadata bson.M `bson:"metadata,omitempty"`

	CreatedAt time.Time `bson:"created_at"`
}
