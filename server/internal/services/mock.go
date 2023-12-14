package services

import (
	"alvinlucillo/swapi-app/internal/models"
	"alvinlucillo/swapi-app/internal/repositories"
)

func NewMockRepository(searches []models.SearchModel, characters []models.CharacterModel, vehicles []models.VehicleModel, films []models.FilmModel) repositories.Repository {
	repository := repositories.Repository{
		SearchRepository: mockSearchRepository{
			searches: searches,
		},
		CharacterRepository: mockCharacterRepository{
			characters: characters,
		},
		VehicleRepository: mockVehicleRepository{
			vehicles: vehicles,
		},
		FilmRepository: mockFilmRepository{
			films: films,
		},
	}

	return repository
}

type MockSWAPIClient struct {
	characters []models.CharacterModel
	vehicles   []models.VehicleModel
	films      []models.FilmModel
}

func NewMockSWAPIClient(searches []models.SearchModel, characters []models.CharacterModel, vehicles []models.VehicleModel, films []models.FilmModel) SWAPIQueryer {
	return MockSWAPIClient{
		characters: characters,
		vehicles:   vehicles,
		films:      films,
	}
}

func (s MockSWAPIClient) QueryPeople(name string) ([]PeopleResult, error) {
	for _, character := range s.characters {
		if character.Name == name {
			return []PeopleResult{
				{
					Name:     character.Name,
					URL:      character.ID,
					Films:    character.Films,
					Vehicles: character.Vehicles,
				},
			}, nil
		}
	}

	return nil, nil
}
func (s MockSWAPIClient) QueryFilm(id string) (FilmResult, error) {
	return FilmResult{
		Title: "A New Hope",
	}, nil
}

func (s MockSWAPIClient) QueryVehicle(id string) (VehicleResult, error) {
	return VehicleResult{
		Model: "T-16 skyhopper",
	}, nil
}

type mockSearchRepository struct {
	searches []models.SearchModel
}

func (m mockSearchRepository) AddSearch(newSearch models.SearchModel) (string, error) {
	m.searches = append(m.searches, newSearch)
	return "", nil
}

func (m mockSearchRepository) GetSearches() ([]models.SearchModel, error) {
	return m.searches, nil
}

func (m mockSearchRepository) GetSearchesByID(searchID string) (*models.SearchModel, error) {
	for _, search := range m.searches {
		if search.ID.Hex() == searchID {
			return &search, nil
		}
	}
	return nil, nil
}

func (m mockSearchRepository) RemoveExpiration(searchID string) (bool, error) {
	for i, search := range m.searches {
		if search.ID.Hex() == searchID {
			m.searches[i].ExpiresAt = nil
			return true, nil
		}
	}
	return false, nil
}

type mockCharacterRepository struct {
	characters []models.CharacterModel
}

func (m mockCharacterRepository) AddCharacter(newCharacter models.CharacterModel) (string, error) {
	m.characters = append(m.characters, newCharacter)
	return "", nil
}

func (m mockCharacterRepository) GetCharacter(id string) (*models.CharacterModel, error) {
	for _, character := range m.characters {
		if character.ID == id {
			return &character, nil
		}
	}
	return nil, nil
}

func (m mockCharacterRepository) GetCharacterByID(id string) (*models.CharacterModel, error) {
	for _, character := range m.characters {
		if character.ID == id {
			return &character, nil
		}
	}
	return nil, nil
}

func (m mockCharacterRepository) GetCharacters() ([]models.CharacterModel, error) {
	return m.characters, nil
}

func (m mockCharacterRepository) GetCharactersByID(ids []string) ([]models.CharacterModel, error) {
	var characters []models.CharacterModel
	for _, id := range ids {
		for _, character := range m.characters {
			if character.ID == id {
				characters = append(characters, character)
			}
		}
	}
	return characters, nil
}

type mockFilmRepository struct {
	films []models.FilmModel
}

func (m mockFilmRepository) AddFilm(newFilm models.FilmModel) (string, error) {
	m.films = append(m.films, newFilm)
	return "", nil
}

func (m mockFilmRepository) GetFilm(id string) (*models.FilmModel, error) {
	for _, film := range m.films {
		if film.ID == id {
			return &film, nil
		}
	}
	return nil, nil
}

type mockVehicleRepository struct {
	vehicles []models.VehicleModel
}

func (m mockVehicleRepository) AddVehicle(newVehicle models.VehicleModel) (string, error) {
	m.vehicles = append(m.vehicles, newVehicle)
	return "", nil
}

func (m mockVehicleRepository) GetVehicle(id string) (*models.VehicleModel, error) {
	for _, vehicle := range m.vehicles {
		if vehicle.ID == id {
			return &vehicle, nil
		}
	}
	return nil, nil
}
