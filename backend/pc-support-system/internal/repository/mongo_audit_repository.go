package repository

import (
	"context"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

type AuditFilter struct {
	Entity      string
	EntityID    *bson.ObjectID
	PerformedBy *bson.ObjectID
	Action      string

	From *time.Time
	To   *time.Time
}

func (r *MongoAuditRepository) FindAuditLogs(
	ctx context.Context,
	filter AuditFilter,
) ([]*models.AuditEvent, error) {

	query := bson.M{}

	if filter.Entity != "" {
		query["entity_type"] = filter.Entity
	}

	if filter.EntityID != nil {
		query["entity_id"] = *filter.EntityID
	}

	if filter.PerformedBy != nil {
		query["performed_by_id"] = *filter.PerformedBy
	}

	if filter.Action != "" {
		query["event_type"] = filter.Action
	}

	if filter.From != nil || filter.To != nil {

		dateFilter := bson.M{}

		if filter.From != nil {
			dateFilter["$gte"] = *filter.From
		}

		if filter.To != nil {
			dateFilter["$lte"] = *filter.To
		}

		query["created_at"] = dateFilter
	}

	cursor, err := r.collection.Find(
		ctx,
		query,
		options.Find().SetSort(
			bson.D{
				{
					Key:   "created_at",
					Value: -1,
				},
			},
		),
	)
	//log.Println("Cursor next:", cursor.Next(ctx))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	logs := make([]*models.AuditEvent, 0)

	for cursor.Next(ctx) {

		var log models.AuditEvent

		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}

		logs = append(logs, &log)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
