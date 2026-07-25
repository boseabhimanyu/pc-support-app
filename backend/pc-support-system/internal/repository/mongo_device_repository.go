package repository

import (
	"context"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoDeviceRepository struct {
	collection *mongo.Collection
}

func NewDeviceRepository(db *mongo.Database) DeviceRepository {
	return &MongoDeviceRepository{
		collection: db.Collection("devices"),
	}
}

func (r *MongoDeviceRepository) Create(
	ctx context.Context,
	device *models.Device,
) error {
	result, err := r.collection.InsertOne(ctx, device)
	if err != nil {
		return err
	}

	device.ID = result.InsertedID.(bson.ObjectID)

	return nil
}

func (r *MongoDeviceRepository) FindByCustomerID(
	ctx context.Context,
	customerID bson.ObjectID,
) ([]models.Device, error) {

	var devices []models.Device

	cursor, err := r.collection.Find(ctx, bson.M{
		"customer_id": customerID,
	})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var device models.Device

		if err := cursor.Decode(&device); err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func (r *MongoDeviceRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (*models.Device, error) {

	var device models.Device

	err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&device)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &device, nil
}
