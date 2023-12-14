package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB is a struct that holds the MongoDB client and database information.
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoDB creates a new instance of MongoDB with the provided connection string, database name, and collection name.
func NewMongoDB(connectionString, dbName string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	database := client.Database(dbName)

	return &MongoDB{
		Client:   client,
		Database: database,
	}, nil
}
