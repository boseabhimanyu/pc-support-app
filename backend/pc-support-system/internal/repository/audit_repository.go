package repository

import (
	"context"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
)

type AuditRepository interface {
	CreateEvent(ctx context.Context, event *models.AuditEvent) error
}
