package repositories

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"alvinlucillo/swapi-app/internal/models"
)

type Repository struct {
	VehicleRepository   VehicleRepository
	FilmRepository      FilmRepository
	SearchRepository    SearchRepository
	CharacterRepository CharacterRepository
}

func NewRepository(db *mongo.Database) (*Repository, error) {
	vehicleRepository, err := NewVehicleRepository(db)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}

	filmRepository, err := NewFilmRepository(db)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}

	searchRepository, err := NewSearchRepository(db)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}

	characterRepository, err := NewCharacterRepository(db)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}

	return &Repository{
		VehicleRepository:   vehicleRepository,
		FilmRepository:      filmRepository,
		SearchRepository:    searchRepository,
		CharacterRepository: characterRepository,
	}, nil
}

type VehicleRepository interface {
	AddVehicle(newVehicle models.VehicleModel) (string, error)
	GetVehicle(id string) (*models.VehicleModel, error)
}

type FilmRepository interface {
	AddFilm(newVehicle models.FilmModel) (string, error)
	GetFilm(id string) (*models.FilmModel, error)
}

type SearchRepository interface {
	AddSearch(newVehicle models.SearchModel) (string, error)
	GetSearches() ([]models.SearchModel, error)
	RemoveExpiration(id string) (bool, error)
	GetSearchesByID(id string) (*models.SearchModel, error)
}

type CharacterRepository interface {
	GetCharacter(id string) (*models.CharacterModel, error)
	AddCharacter(newCharacter models.CharacterModel) (string, error)
}
