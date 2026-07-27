package repository

import (
	"context"
	"strings"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) Create(ctx context.Context, user *models.User) error {

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(bson.ObjectID)

	return nil
}

func (r *MongoUserRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User

	err := r.collection.FindOne(ctx, bson.M{
		"phone": phone,
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 👈 IMPORTANT: not found is NOT an error
		}
		return nil, err // real DB error
	}

	return &user, nil // user found
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.collection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id bson.ObjectID) (*models.User, error) {
	var user models.User

	err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *MongoUserRepository) Update(
	ctx context.Context,
	user *models.User,
) error {

	update := bson.M{
		"$set": bson.M{
			"first_name":    user.FirstName,
			"last_name":     user.LastName,
			"email":         user.Email,
			"phone":         user.Phone,
			"updated_at":    user.UpdatedAt,
			"updated_by_id": user.UpdatedByID,
		},
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		update,
	)

	return err
}

func (r *MongoUserRepository) Search(
	ctx context.Context,
	query string,
) ([]models.User, error) {

	filter := bson.M{
		"role": models.RoleCustomer,
		"$or": []bson.M{
			{
				"phone": bson.M{
					"$regex":   query,
					"$options": "i",
				},
			},
			{
				"email": bson.M{
					"$regex":   query,
					"$options": "i",
				},
			},
			{
				"first_name": bson.M{
					"$regex":   query,
					"$options": "i",
				},
			},
			{
				"last_name": bson.M{
					"$regex":   query,
					"$options": "i"},
			},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User

	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *MongoUserRepository) UpdatePassword(
	ctx context.Context,
	user *models.User,
) error {

	update := bson.M{
		"$set": bson.M{
			"password_hash": user.PasswordHash,
			"updated_at":    user.UpdatedAt,
			"updated_by_id": user.UpdatedByID,
		},
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		update,
	)

	return err
}

func (r *MongoUserRepository) UpdateCustomer(
	ctx context.Context,
	user *models.User,
) error {

	filter := bson.M{
		"_id": user.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"first_name":    user.FirstName,
			"last_name":     user.LastName,
			"phone":         user.Phone,
			"email":         user.Email,
			"updated_at":    user.UpdatedAt,
			"updated_by_id": user.UpdatedByID,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoUserRepository) SearchStaff(
	ctx context.Context,
	query string,
) ([]models.User, error) {

	query = strings.TrimSpace(query)

	filter := bson.M{
		"role": bson.M{
			"$ne": models.RoleCustomer,
		},
	}

	if query != "" {

		regex := bson.Regex{
			Pattern: query,
			Options: "i",
		}

		filter = bson.M{
			"role": bson.M{
				"$ne": models.RoleCustomer,
			},
			"$or": []bson.M{
				{"first_name": regex},
				{"last_name": regex},
				{"email": regex},
				{"phone": regex},
				{"role": regex},
			},
		}
	}

	opts := options.Find().
		SetSort(bson.D{
			{Key: "first_name", Value: 1},
			{Key: "last_name", Value: 1},
		})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var staff []models.User

	if err := cursor.All(ctx, &staff); err != nil {
		return nil, err
	}

	return staff, nil
}
