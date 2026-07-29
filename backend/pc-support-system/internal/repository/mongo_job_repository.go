package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (r *MongoJobRepository) GenerateJobNumber(
	ctx context.Context,
) (string, error) {

	today := time.Now().Format("20060102")

	prefix := "JOB-" + today + "-"

	filter := bson.M{
		"job_number": bson.M{
			"$regex": "^" + prefix,
		},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return "", err
	}

	next := count + 1

	jobNumber := fmt.Sprintf("%s%04d", prefix, next)

	return jobNumber, nil
}

func (r *MongoJobRepository) FindOpenUnassigned(
	ctx context.Context,
) ([]*models.Job, error) {

	filter := bson.M{
		"status":         models.JobCreated,
		"assigned_to_id": bson.M{"$exists": false},
	}

	opts := options.Find().
		SetSort(bson.D{
			{Key: "created_at", Value: 1},
		})

	cursor, err := r.collection.Find(
		ctx,
		filter,
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []*models.Job

	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *MongoJobRepository) AssignJob(
	ctx context.Context,
	job *models.Job,
) error {

	filter := bson.M{
		"_id": job.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"assigned_to_id": job.AssignedToID,
			"assigned_by_id": job.AssignedByID,
			"status":         job.Status,
			"updated_at":     job.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}

func (r *MongoJobRepository) FindOpenAssigned(
	ctx context.Context,
) ([]*models.Job, error) {

	filter := bson.M{
		"status":         models.JobAssigned,
		"assigned_to_id": bson.M{"$exists": true},
	}

	opts := options.Find().
		SetSort(bson.D{
			{Key: "created_at", Value: 1},
		})

	cursor, err := r.collection.Find(
		ctx,
		filter,
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []*models.Job

	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *MongoJobRepository) FindMyJobs(
	ctx context.Context,
	staffID bson.ObjectID,
) ([]*models.Job, error) {

	filter := bson.M{
		"assigned_to_id": staffID,
		"status": bson.M{
			"$in": []models.JobStatus{
				models.JobAssigned,
				models.JobInProgress,
				models.JobWaitingCustomer,
				models.JobResumed,
			},
		},
	}

	opts := options.Find().
		SetSort(bson.D{
			{Key: "updated_at", Value: -1},
		})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []*models.Job

	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *MongoJobRepository) FindCustomerJobs(
	ctx context.Context,
	customerID bson.ObjectID,
) ([]*models.Job, error) {

	filter := bson.M{
		"customer_id": customerID,
	}

	opts := options.Find().
		SetSort(bson.D{
			{Key: "created_at", Value: -1},
		})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []*models.Job

	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *MongoJobRepository) UpdateStatus(
	ctx context.Context,
	jobID bson.ObjectID,
	status models.JobStatus,
) error {

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateByID(
		ctx,
		jobID,
		update,
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *MongoJobRepository) FindByStatuses(
	ctx context.Context,
	statuses []models.JobStatus,
) ([]*models.Job, error) {

	filter := bson.M{
		"status": bson.M{
			"$in": statuses,
		},
	}

	opts := options.Find().
		SetSort(bson.D{
			{Key: "created_at", Value: 1},
		})

	cursor, err := r.collection.Find(
		ctx,
		filter,
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []*models.Job

	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}
