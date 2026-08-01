package repository

import (
	"context"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	FindByID(ctx context.Context, id bson.ObjectID) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Search(ctx context.Context, query string) ([]models.User, error)
	UpdatePassword(ctx context.Context, user *models.User) error
	UpdateCustomer(ctx context.Context, user *models.User) error
	SearchStaff(ctx context.Context, query string) ([]models.User, error)
	UpdateRefreshToken(ctx context.Context, userID bson.ObjectID, token string, expiresAt time.Time) error
	ClearRefreshToken(ctx context.Context, userID bson.ObjectID) error
	FindByRefreshToken(ctx context.Context, token string) (*models.User, error)
	// Delete(id bson.ObjectID) error
}
