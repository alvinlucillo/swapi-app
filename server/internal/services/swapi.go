package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SWAPI is a client for the Star Wars API
type SWAPIClient struct {
	client  *http.Client
	baseURL string // https://swapi.dev/api
}

// SWAPIQueryer is an interface for querying the Star Wars API
type SWAPIQueryer interface {
	QueryPeople(name string) ([]PeopleResult, error)
	QueryFilm(filmID string) (FilmResult, error)
	QueryVehicle(vehicleID string) (VehicleResult, error)
}

type PeopleResult struct {
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Vehicles []string `json:"vehicles"`
	Films    []string `json:"films"`
}

type VehicleResult struct {
	Model string `json:"model"`
	URL   string `json:"url"`
}

type FilmResult struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func NewSWAPIClient(client *http.Client, baseURL string) SWAPIClient {
	return SWAPIClient{client: client, baseURL: baseURL}
}

type PeopleResponse struct {
	Results []PeopleResult `json:"results"`
}

type FilmResponse struct {
	Result FilmResult `json:"results"`
}

// QueryPeople - queries the Star Wars API for people with the given name
func (s SWAPIClient) QueryPeople(name string) ([]PeopleResult, error) {
	escapedName := url.QueryEscape(name)
	url := s.baseURL + "/people/?search=" + escapedName

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response PeopleResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Results, nil
}

// QueryFilm - queries the Star Wars API for a film with the given ID
func (s SWAPIClient) QueryFilm(sourceUrl string) (FilmResult, error) {
	var response FilmResult

	req, err := http.NewRequest(http.MethodGet, sourceUrl, nil)
	if err != nil {
		return response, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return response, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

// QueryVehicle - queries the Star Wars API for a vehicle with the given ID
func (s SWAPIClient) QueryVehicle(sourceUrl string) (VehicleResult, error) {
	req, err := http.NewRequest(http.MethodGet, sourceUrl, nil)
	if err != nil {
		return VehicleResult{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return VehicleResult{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var result VehicleResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return VehicleResult{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
