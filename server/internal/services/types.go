package services

type Character struct {
	ID            string
	Name          string
	VehicleModels []string
	Films         []string
}

type Search struct {
	ID        string
	SearchKey string
}
