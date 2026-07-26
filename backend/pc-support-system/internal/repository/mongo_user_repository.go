package repository

import (
	"context"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"phone":      user.Phone,
			"updated_at": user.UpdatedAt,
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
	id bson.ObjectID,
	passwordHash string,
) error {

	update := bson.M{
		"$set": bson.M{
			"password_hash": passwordHash,
			"updated_at":    time.Now(),
		},
	}

	_, err := r.collection.UpdateByID(ctx, id, update)

	return err
}
