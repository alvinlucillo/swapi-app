package internal

import "fmt"

type CharacterService interface {
	GetCharacters(name string) ([]Character, error)
	GetSavedSearches() ([]Character, error)
}

type CharacterServiceImpl struct {
	swapiClient SWAPIQueryer
}

func NewService(swapiClient SWAPIQueryer) CharacterServiceImpl {
	return CharacterServiceImpl{
		swapiClient: swapiClient,
	}
}

func (c CharacterServiceImpl) GetCharacters(name string) ([]Character, error) {
	peopleResult, err := c.swapiClient.QueryPeople(name)
	if err != nil {
		return nil, fmt.Errorf("failed to query people: %w", err)
	}

	var characters []Character
	for _, person := range peopleResult {
		character := Character{
			ID:           person.URL,
			Name:         person.Name,
			IsBookmarked: false,
		}

		for _, film := range person.Films {
			filmResult, err := c.swapiClient.QueryFilm(film)
			if err != nil {
				return nil, fmt.Errorf("failed to query film: %w", err)
			}

			character.Films = append(character.Films, filmResult.Title)
		}

		for _, vehicle := range person.Vehicles {
			vehicleResult, err := c.swapiClient.QueryVehicle(vehicle)
			if err != nil {
				return nil, fmt.Errorf("failed to query vehicle: %w", err)
			}

			character.CarModels = append(character.CarModels, vehicleResult.Model)
		}

		characters = append(characters, character)

	}

	return characters, nil
}

func (CharacterServiceImpl) GetSavedSearches() ([]Character, error) {

	Luke := Character{
		ID:        "1000",
		Name:      "Luke Skywalker",
		CarModels: []string{"X-34 Landspeeder", "T-16 Skyhopper"},
		Films:     []string{"A New Hope", "The Empire Strikes Back", "Return of the Jedi"},
	}

	return []Character{Luke}, nil
}
