package repository

import (
	"context"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobRepository interface {
	Create(ctx context.Context, job *models.Job) error

	FindOpenUnassigned(ctx context.Context) ([]*models.Job, error)

	FindOpenAssigned(ctx context.Context) ([]*models.Job, error)

	FindMyJobs(ctx context.Context, staffID bson.ObjectID) ([]*models.Job, error)

	AssignJob(ctx context.Context, job *models.Job) error

	FindByID(ctx context.Context, id bson.ObjectID) (*models.Job, error)

	FindCustomerJobs(ctx context.Context, customerID bson.ObjectID) ([]*models.Job, error)

	UpdateStatus(ctx context.Context, jobID bson.ObjectID, status models.JobStatus) error

	FindByJobNumber(ctx context.Context, jobNumber string) (*models.Job, error)

	FindByStatuses(ctx context.Context, statuses []models.JobStatus) ([]*models.Job, error)

	AddNote(ctx context.Context, jobID bson.ObjectID, note models.JobNote) error

	CloseJob(ctx context.Context, job *models.Job) error

	SearchJobs(ctx context.Context, query string) ([]*models.Job, error)

	//FindOpenJobByDeviceID(ctx, deviceID) // To find anyexisting open job for the same device

	GenerateJobNumber(ctx context.Context) (string, error)
}
