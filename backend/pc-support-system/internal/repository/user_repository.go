package repository

import (
	"context"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	// FindByEmail(ctx context.Context, email string) (*models.User, error)
	// FindByPhone(ctx context.Context, phone string) (*models.User, error)
	// FindByID(ctx context.Context, id bson.ObjectID) (*models.User, error)
}

// Update(user *models.User) error

// Delete(id bson.ObjectID) error
