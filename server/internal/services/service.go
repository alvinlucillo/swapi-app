package services

import (
	"fmt"

	"alvinlucillo/swapi-app/internal/db"
	"alvinlucillo/swapi-app/internal/models"
	"alvinlucillo/swapi-app/internal/repositories"

	"github.com/caarlos0/env/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	DBName             string `env:"DB_NAME" envDefault:"swapiapp"`
	DBConnectionString string `env:"DB_CONNECTION_STRING" envDefault:"mongodb://localhost:27017/"`
	DBDocumentTTL      int32  `env:"DB_DOCUMENT_TTL" envDefault:"43200"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

type CharacterService interface {
	GetCharacters(name string) ([]Character, string, error)
	GetSavedSearches() ([]Search, error)
	GetSavedSearchesByID(searchID string) ([]Character, error)
	SaveSearch(searchID string) (bool, error)
}

type CharacterServiceImpl struct {
	swapiClient SWAPIQueryer
	repository  *repositories.Repository
}

func NewService(swapiClient SWAPIQueryer) (*CharacterServiceImpl, error) {

	cfg, err := NewConfig()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}

	db, err := db.NewMongoDB(cfg.DBConnectionString, cfg.DBName)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}

	repo, err := repositories.NewRepository(repositories.Config{DocumentTTL: cfg.DBDocumentTTL, DB: db.Database})
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}
	return &CharacterServiceImpl{
		swapiClient: swapiClient,
		repository:  repo,
	}, nil
}

// GetCharacters -
//  1. Queries the SWAPI for people with the given name
//  2. Adds the films and vehicles to the database if they don't already exist
//  3. Adds the character to the database if it doesn't already exist
//  4. Adds the search to the database for retrieval later
func (c CharacterServiceImpl) GetCharacters(name string) ([]Character, string, error) {
	peopleResult, err := c.swapiClient.QueryPeople(name)
	if err != nil {
		return nil, "", fmt.Errorf("failed to query people: %w", err)
	}

	var characters []Character
	var characterIDs, filmIDs, vehicleIDs []string
	for _, person := range peopleResult {
		character := Character{
			ID:   person.URL,
			Name: person.Name,
		}

		for _, film := range person.Films {
			existingFilm, err := c.repository.FilmRepository.GetFilm(film)
			if err != nil {
				return nil, "", fmt.Errorf("failed to get film: %w", err)
			}

			var f models.FilmModel
			if existingFilm != nil {
				f = *existingFilm
			} else {
				filmResult, err := c.swapiClient.QueryFilm(film)
				if err != nil {
					return nil, "", fmt.Errorf("failed to query film: %w", err)
				}

				f = models.FilmModel{
					Title: filmResult.Title,
					ID:    filmResult.URL,
				}
				_, err = c.repository.FilmRepository.AddFilm(f)
				if err != nil {
					return nil, "", fmt.Errorf("failed to add film: %w", err)
				}
			}
			filmIDs = append(filmIDs, film)
			character.Films = append(character.Films, f.Title)
		}

		for _, vehicle := range person.Vehicles {
			existingVehicle, err := c.repository.VehicleRepository.GetVehicle(vehicle)
			if err != nil {
				return nil, "", fmt.Errorf("failed to get vehicle: %w", err)
			}

			var v models.VehicleModel
			if existingVehicle != nil {
				v = *existingVehicle
			} else {
				vehicleResult, err := c.swapiClient.QueryVehicle(vehicle)
				if err != nil {
					return nil, "", fmt.Errorf("failed to query vehicle: %w", err)
				}

				v = models.VehicleModel{
					Model: vehicleResult.Model,
					ID:    vehicleResult.URL,
				}
				_, err = c.repository.VehicleRepository.AddVehicle(v)
				if err != nil {
					return nil, "", fmt.Errorf("failed to add vehicle: %w", err)
				}
			}
			vehicleIDs = append(vehicleIDs, vehicle)
			character.VehicleModels = append(character.VehicleModels, v.Model)
		}

		characterIDs = append(characterIDs, character.ID)
		characters = append(characters, character)

		existingCharacter, error := c.repository.CharacterRepository.GetCharacter(person.URL)
		if error != nil {
			return nil, "", fmt.Errorf("failed to get character: %w", err)
		}

		if existingCharacter == nil {
			_, error := c.repository.CharacterRepository.AddCharacter(models.CharacterModel{
				Name:     person.Name,
				ID:       person.URL,
				Films:    filmIDs,
				Vehicles: vehicleIDs,
			})
			if error != nil {
				return nil, "", fmt.Errorf("failed to add character: %w", err)
			}
		}
	}

	if len(peopleResult) == 0 {
		return nil, "", nil
	}

	search := models.SearchModel{
		ID:         primitive.NewObjectID(),
		SearchKey:  name,
		Characters: characterIDs,
	}

	searchID, err := c.repository.SearchRepository.AddSearch(search)
	if err != nil {
		return nil, "", fmt.Errorf("failed to add search: %w", err)
	}

	return characters, searchID, nil
}

// GetSavedSearches - Gets saved searches from the database
func (c CharacterServiceImpl) GetSavedSearches() ([]Search, error) {
	searches, err := c.repository.SearchRepository.GetSearches()
	if err != nil {
		return nil, fmt.Errorf("failed to get searches: %w", err)
	}

	var searchResults []Search
	for _, search := range searches {
		searchResults = append(searchResults, Search{
			ID:        search.ID.Hex(),
			SearchKey: search.SearchKey,
		})
	}

	return searchResults, nil
}

// SaveSearch - Removes the expiration from a search so it's not marked for deletion
func (c CharacterServiceImpl) SaveSearch(searchID string) (bool, error) {
	result, err := c.repository.SearchRepository.RemoveExpiration(searchID)
	if err != nil {
		return false, fmt.Errorf("failed to remove expiration: %w", err)
	}

	return result, nil
}

// GetSavedSearchesByID - Gets saved searches from the database by ID
// 1. Gets the characters from the search
// 2. Builds the character result with the film and vehicle data
func (c CharacterServiceImpl) GetSavedSearchesByID(searchID string) ([]Character, error) {
	search, err := c.repository.SearchRepository.GetSearchesByID(searchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get search by ID: %w", err)
	}

	if search == nil {
		return nil, nil
	}

	var characters []Character
	for _, characterID := range search.Characters {
		character, err := c.repository.CharacterRepository.GetCharacter(characterID)
		if err != nil {
			return nil, fmt.Errorf("failed to get character: %w", err)
		}

		characterResult := Character{
			ID:   character.ID,
			Name: character.Name,
		}

		for _, film := range character.Films {
			existingFilm, err := c.repository.FilmRepository.GetFilm(film)
			if err != nil {
				return nil, fmt.Errorf("failed to get film: %w", err)
			}

			characterResult.Films = append(characterResult.Films, existingFilm.Title)
		}

		for _, vehicle := range character.Vehicles {
			existingVehicle, err := c.repository.VehicleRepository.GetVehicle(vehicle)
			if err != nil {
				return nil, fmt.Errorf("failed to get vehicle: %w", err)
			}

			characterResult.VehicleModels = append(characterResult.VehicleModels, existingVehicle.Model)
		}

		characters = append(characters, characterResult)
	}

	return characters, nil
}
