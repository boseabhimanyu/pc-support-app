package repository

import (
	"context"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoAuditRepository struct {
	collection *mongo.Collection
}

func NewAuditRepository(db *mongo.Database) AuditRepository {
	return &MongoAuditRepository{
		collection: db.Collection("audit_events"),
	}
}

func (r *MongoAuditRepository) CreateEvent(
	ctx context.Context,
	event *models.AuditEvent,
) error {

	_, err := r.collection.InsertOne(
		ctx,
		event,
	)

	return err
}
