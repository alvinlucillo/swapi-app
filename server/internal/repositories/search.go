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
	SearchCollection = "searches"
)

type SearchRepositoryImpl struct {
	db *mongo.Database
}

func NewSearchRepository(cfg Config) (*SearchRepositoryImpl, error) {

	collection := cfg.DB.Collection(SearchCollection)

	// Define the index model
	// Does not automatically expire the documents; expiration is based on the value of the expiresAt field
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	cursor, err := collection.Indexes().List(context.TODO())
	if err != nil {
		return nil, err
	}
	var indexes []bson.M
	if err = cursor.All(context.TODO(), &indexes); err != nil {
		return nil, err
	}

	// Create the index if it doesn't exist
	createIndex := true
	for _, index := range indexes {
		if index["name"].(string) == "expiresAt_1" {
			createIndex = false
			break
		}
	}

	if createIndex {
		// create the index
		_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
		if err != nil {
			return nil, err
		}
	}

	return &SearchRepositoryImpl{
		db: cfg.DB,
	}, nil
}

// AddSearch - Adds a new search to the database
func (r *SearchRepositoryImpl) AddSearch(search models.SearchModel) (string, error) {
	collection := r.db.Collection(SearchCollection)

	// Document by default will expire 1 hour after creation
	t := time.Now().Add(1 * time.Hour)
	search.ExpiresAt = &t

	result, err := collection.InsertOne(context.TODO(), search)
	if err != nil {
		fmt.Printf("failed to insert search: %v", err)
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetSearches - Returns all searches that have not expired
func (r *SearchRepositoryImpl) GetSearches() ([]models.SearchModel, error) {
	collection := r.db.Collection(SearchCollection)

	// Find all documents with nil expiresAt
	filter := bson.D{{Key: "expiresAt", Value: nil}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var searches []models.SearchModel
	if err = cursor.All(context.Background(), &searches); err != nil {
		return nil, err
	}

	return searches, nil
}

// RemoveExpiration - Removes the expiration from a search
func (r *SearchRepositoryImpl) RemoveExpiration(searchID string) (bool, error) {
	collection := r.db.Collection(SearchCollection)

	// Convert searchID to ObjectId
	objectID, err := primitive.ObjectIDFromHex(searchID)
	if err != nil {
		return false, err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}

	// Set the expiresAt field to nil
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "expiresAt", Value: nil}}}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false, err
	}

	return result.ModifiedCount > 0, nil
}

// GetSearchByID - Returns a search by ID
func (r *SearchRepositoryImpl) GetSearchesByID(searchID string) (*models.SearchModel, error) {
	collection := r.db.Collection(SearchCollection)

	// Convert searchID to ObjectId
	objectID, err := primitive.ObjectIDFromHex(searchID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}

	var search models.SearchModel
	err = collection.FindOne(context.Background(), filter).Decode(&search)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle no document found
			return nil, nil
		}
		return nil, err
	}

	return &search, nil
}
