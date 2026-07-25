package repository

import (
	"context"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoJobRepository struct {
	collection *mongo.Collection
}

func NewJobRepository(db *mongo.Database) JobRepository {
	return &MongoJobRepository{
		collection: db.Collection("jobs"),
	}
}

func (r *MongoJobRepository) Create(
	ctx context.Context,
	job *models.Job,
) error {

	result, err := r.collection.InsertOne(ctx, job)
	if err != nil {
		return err
	}

	job.ID = result.InsertedID.(bson.ObjectID)

	return nil
}

func (r *MongoJobRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (*models.Job, error) {

	var job models.Job

	err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&job)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *MongoJobRepository) FindByJobNumber(
	ctx context.Context,
	jobNumber string,
) (*models.Job, error) {

	var job models.Job

	err := r.collection.FindOne(ctx, bson.M{
		"job_number": jobNumber,
	}).Decode(&job)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &job, nil
}

// func (r *MongoJobRepository) FindByStatus(
// 	ctx context.Context,
// 	status models.JobStatus,
// ) ([]models.Job, error)
