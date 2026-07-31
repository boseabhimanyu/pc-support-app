package services

import (
	"context"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuditService struct {
	repo repository.AuditRepository
}

func NewAuditService(
	repo repository.AuditRepository,
) *AuditService {

	return &AuditService{
		repo: repo,
	}
}

func (s *AuditService) Log(
	ctx context.Context,
	entityType models.AuditEntityType,
	entityID bson.ObjectID,
	eventType models.AuditEventType,
	performedBy bson.ObjectID,
	metadata bson.M,
) error {

	event := &models.AuditEvent{
		EntityType:    entityType,
		EntityID:      entityID,
		EventType:     eventType,
		PerformedByID: performedBy,
		Metadata:      metadata,
		CreatedAt:     time.Now(),
	}

	return s.repo.CreateEvent(
		ctx,
		event,
	)
}
