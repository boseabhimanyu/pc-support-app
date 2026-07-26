package repository

import (
	"context"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobRepository interface {
	Create(ctx context.Context, job *models.Job) error

	//FindOpenJobByDeviceID(ctx, deviceID) // To find anyexisting open job for the same device

	FindByID(ctx context.Context, id bson.ObjectID) (*models.Job, error)

	FindByJobNumber(ctx context.Context, jobNumber string) (*models.Job, error)

	//FindByStatus(ctx context.Context, status models.JobStatus) ([]models.Job, error)

	GenerateJobNumber(ctx context.Context) (string, error)
}
