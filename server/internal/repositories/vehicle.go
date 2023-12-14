package repositories

import (
	"context"
	"fmt"
	"time"

	"alvinlucillo/swapi-app/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	VehicleCollection = "vehicles"
)

type VehicleRepositoryImpl struct {
	db *mongo.Database
}

func NewVehicleRepository(cfg Config) (*VehicleRepositoryImpl, error) {

	collection := cfg.DB.Collection(VehicleCollection)

	// Define the index model
	// Sets TTL
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "createdAt", Value: 1}},                   // Index key
		Options: options.Index().SetExpireAfterSeconds(cfg.DocumentTTL), // TTL value
	}

	cursor, err := collection.Indexes().List(context.TODO())
	if err != nil {
		return nil, err
	}
	var indexes []bson.M
	if err = cursor.All(context.TODO(), &indexes); err != nil {
		return nil, err
	}

	// Create a map of index names to expireAfterSeconds values
	indexMap := map[string]int32{}
	for _, index := range indexes {
		if index["expireAfterSeconds"] != nil {
			indexMap[index["name"].(string)] = index["expireAfterSeconds"].(int32)
		}
	}

	// Create index if it doesn't exist or if it exists but has a different expireAfterSeconds value
	createIndex := false
	if expireSec := indexMap["createdAt_1"]; expireSec != 0 {
		if expireSec != 3600 {
			// index exists but has a different expireAfterSeconds value so drop it and create a new one
			_, err := collection.Indexes().DropOne(context.TODO(), "createdAt_1")
			if err != nil {
				return nil, err
			}
			createIndex = true
		} else {
			// nothing to do; index already exists and has the same expireAfterSeconds value
		}
	} else {
		// index does not exist so create it
		createIndex = true
	}

	if createIndex {
		// create the index
		_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
		if err != nil {
			return nil, err
		}
	}

	return &VehicleRepositoryImpl{
		db: cfg.DB,
	}, nil
}

// AddVehicle - adds a vehicle to the database
func (r *VehicleRepositoryImpl) AddVehicle(vehicle models.VehicleModel) (string, error) {
	collection := r.db.Collection(VehicleCollection)

	vehicle.CreatedAt = time.Now()

	result, err := collection.InsertOne(context.TODO(), vehicle)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetVehicle - returns a vehicle from the database
func (r *VehicleRepositoryImpl) GetVehicle(url string) (*models.VehicleModel, error) {
	collection := r.db.Collection(VehicleCollection)

	filter := bson.M{"id": url}

	var vehicle models.VehicleModel
	err := collection.FindOne(context.TODO(), filter).Decode(&vehicle)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle no document found
			fmt.Println("no document matches the provided filter for vehicle")
			return nil, nil
		}
		return nil, err
	}

	return &vehicle, nil
}
