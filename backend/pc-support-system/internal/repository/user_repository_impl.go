package repository

import (
	"context"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *MongoUserRepository) UpdateRefreshToken(
	ctx context.Context,
	userID bson.ObjectID,
	token string,
	expiresAt time.Time,
) error {

	update := bson.M{
		"$set": bson.M{
			"current_refresh_token":    token,
			"refresh_token_expires_at": expiresAt,
		},
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		update,
	)

	return err
}

func (r *MongoUserRepository) FindByRefreshToken(
	ctx context.Context,
	token string,
) (*models.User, error) {

	var user models.User

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"current_refresh_token": token,
		},
	).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *MongoUserRepository) ClearRefreshToken(
	ctx context.Context,
	userID bson.ObjectID,
) error {

	update := bson.M{
		"$unset": bson.M{
			"current_refresh_token":    "",
			"refresh_token_expires_at": "",
		},
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		update,
	)

	return err
}
