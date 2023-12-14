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
	FilmCollection = "films"
)

type FilmRepositoryImpl struct {
	db *mongo.Database
}

func NewFilmRepository(db *mongo.Database) (*FilmRepositoryImpl, error) {

	collection := db.Collection(FilmCollection)

	// Define the index model
	// Sets TTL
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "createdAt", Value: 1}},        // Index key
		Options: options.Index().SetExpireAfterSeconds(3600), // TTL value - 1 hour
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

	return &FilmRepositoryImpl{
		db: db,
	}, nil
}

// AddFilm - Adds a film to the database
func (r *FilmRepositoryImpl) AddFilm(film models.FilmModel) (string, error) {
	collection := r.db.Collection(FilmCollection)

	film.CreatedAt = time.Now()

	result, err := collection.InsertOne(context.TODO(), film)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetFilm - Gets a film from the database
func (r *FilmRepositoryImpl) GetFilm(id string) (*models.FilmModel, error) {
	collection := r.db.Collection(FilmCollection)

	filter := bson.M{"id": id}

	var film models.FilmModel
	err := collection.FindOne(context.TODO(), filter).Decode(&film)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle no document found
			fmt.Println("no document matches the provided filter for film")
			return nil, nil
		}
		return nil, err
	}

	return &film, nil
}
