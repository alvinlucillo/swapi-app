package repositories

import (
	"fmt"

	"alvinlucillo/swapi-app/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DocumentTTL int32
)

type Repository struct {
	VehicleRepository   VehicleRepository
	FilmRepository      FilmRepository
	SearchRepository    SearchRepository
	CharacterRepository CharacterRepository
}

type Config struct {
	DocumentTTL int32
	DB          *mongo.Database
}

func NewRepository(cfg Config) (*Repository, error) {
	vehicleRepository, err := NewVehicleRepository(cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}

	filmRepository, err := NewFilmRepository(cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}

	searchRepository, err := NewSearchRepository(cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}

	characterRepository, err := NewCharacterRepository(cfg)
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
