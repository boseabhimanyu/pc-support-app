package repository

import (
	"context"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type DeviceRepository interface {
	Create(ctx context.Context, device *models.Device) error

	FindByID(
		ctx context.Context,
		id bson.ObjectID,
	) (*models.Device, error)

	FindByCustomerID(
		ctx context.Context,
		customerID bson.ObjectID,
	) ([]models.Device, error)

	// Update(ctx context.Context, device *models.Device) error
}
