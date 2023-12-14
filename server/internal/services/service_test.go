package services

import (
	"alvinlucillo/swapi-app/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func generateMockData() ([]models.SearchModel, []models.CharacterModel, []models.VehicleModel, []models.FilmModel) {
	characters := []models.CharacterModel{
		{ID: "1", Name: "Luke Skywalker", Films: []string{"1", "2", "3"}, Vehicles: []string{"1", "2", "3"}},
	}

	vehicles := []models.VehicleModel{
		{ID: "1", Model: "T-16 skyhopper"},
		{ID: "2", Model: "X-34 landspeeder"},
		{ID: "3", Model: "TIE/LN starfighter"},
	}

	films := []models.FilmModel{
		{ID: "1", Title: "A New Hope"},
		{ID: "2", Title: "The Empire Strikes Back"},
		{ID: "3", Title: "Return of the Jedi"},
	}

	timeNow := time.Now()
	searches := []models.SearchModel{
		{ID: primitive.NewObjectID(), SearchKey: "Luke Skywalker", Characters: []string{"1"}, ExpiresAt: nil},
		{ID: primitive.NewObjectID(), SearchKey: "Darth Vader", Characters: []string{"2"}, ExpiresAt: &timeNow},
	}

	return searches, characters, vehicles, films
}

func TestGetSavedSearches(t *testing.T) {
	searches, characters, vehicles, films := generateMockData()
	repository := NewMockRepository(searches, characters, vehicles, films)

	svc := CharacterServiceImpl{
		repository: &repository,
	}

	searchResult, err := svc.GetSavedSearches()
	require.NoError(t, err, "error should be nil")

	require.Equal(t, len(searches), len(searchResult), "searches should be equal")
	require.Equal(t, searches[0].ID.Hex(), searchResult[0].ID, "searches should be equal")
	require.Equal(t, searches[0].SearchKey, searchResult[0].SearchKey, "searches should be equal")
}

func TestSaveSearch(t *testing.T) {
	searches, characters, vehicles, films := generateMockData()
	repository := NewMockRepository(searches, characters, vehicles, films)

	svc := CharacterServiceImpl{
		repository: &repository,
	}

	searchResult, err := svc.SaveSearch(searches[1].ID.Hex())
	require.NoError(t, err, "error should be nil")

	require.Equal(t, true, searchResult, "result should be equal")
}

func TestGetSavedSearchesByID(t *testing.T) {
	searches, characters, vehicles, films := generateMockData()
	repository := NewMockRepository(searches, characters, vehicles, films)

	svc := CharacterServiceImpl{
		repository: &repository,
	}

	searchResult, err := svc.GetSavedSearchesByID(searches[0].ID.Hex())
	require.NoError(t, err, "error should be nil")

	require.Equal(t, len(searchResult), 1, "character length should be equal")
	require.Equal(t, searchResult[0].ID, characters[0].ID, "character ID should be equal")
	require.Equal(t, searchResult[0].Name, characters[0].Name, "character Name should be equal")
	require.Equal(t, len(searchResult[0].Films), 3, "film length should be equal")
	require.Equal(t, len(searchResult[0].VehicleModels), 3, "vehicle length should be equal")
}

func TestGetCharacter(t *testing.T) {
	searches, characters, vehicles, films := generateMockData()
	repository := NewMockRepository(searches, characters, vehicles, films)
	swapiClient := NewMockSWAPIClient(searches, characters, vehicles, films)

	svc := CharacterServiceImpl{
		repository:  &repository,
		swapiClient: swapiClient,
	}

	searchResult, _, err := svc.GetCharacters(characters[0].Name)
	require.NoError(t, err, "error should be nil")

	require.Equal(t, len(searchResult), 1, "character length should be equal")
	require.Equal(t, searchResult[0].ID, characters[0].ID, "character ID should be equal")
	require.Equal(t, searchResult[0].Name, characters[0].Name, "character Name should be equal")
	require.Equal(t, len(searchResult[0].Films), 3, "film length should be equal")
	require.Equal(t, len(searchResult[0].VehicleModels), 3, "vehicle length should be equal")

}
