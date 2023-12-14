package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchModel struct {
	ID         primitive.ObjectID `bson:"_id"`
	SearchKey  string             `bson:"searchKey"`
	Characters []string           `bson:"characters"`
	ExpiresAt  *time.Time         `bson:"expiresAt"`
}

type CharacterModel struct {
	ID        string    `bson:"id"`
	Name      string    `bson:"name"`
	Films     []string  `bson:"films"`
	Vehicles  []string  `bson:"vehicles"`
	CreatedAt time.Time `bson:"createdAt"`
}

type FilmModel struct {
	ID        string    `bson:"id"`
	Title     string    `bson:"title"`
	CreatedAt time.Time `bson:"createdAt"`
}

type VehicleModel struct {
	ID        string    `bson:"id"`
	Model     string    `bson:"model"`
	CreatedAt time.Time `bson:"createdAt"`
}
