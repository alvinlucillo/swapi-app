package services

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type HandlerConfig struct {
	Pretty   bool
	GraphiQL bool
}

type CharactersResult struct {
	Characters []Character
	SearchID   string
}

// NewHandler returns a new graphql handler
func NewHandler(cfg HandlerConfig, svc CharacterService) *handler.Handler {
	// Defines the properties of a character
	// For example, R2-D2 has the following properties:
	// {
	//   "name": "R2-D2",
	//   "films": [
	//     "https://swapi.dev/api/films/1/",
	//   ],
	//   "vehicleModels": [
	//     "https://swapi.dev/api/vehicles/8/"
	//   ]
	// }
	characterType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Character",
		Description: "A StarWars character",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "The name of the character.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if character, ok := p.Source.(Character); ok {
						return character.Name, nil
					}
					return nil, nil
				},
			},
			"films": &graphql.Field{
				Type:        graphql.NewList(graphql.String),
				Description: "The films the character has been in.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if character, ok := p.Source.(Character); ok {
						return character.Films, nil
					}
					return []interface{}{}, nil
				},
			},
			"vehicleModels": &graphql.Field{
				Type:        graphql.NewList(graphql.String),
				Description: "The vehicle models the character drives.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if character, ok := p.Source.(Character); ok {
						return character.VehicleModels, nil
					}
					return []interface{}{}, nil
				},
			},
		},
	})

	// Defines the list of characters that are returned from a search with its search ID
	charactersResultType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "CharactersResult",
			Fields: graphql.Fields{
				"Characters": &graphql.Field{
					Type: graphql.NewList(characterType),
				},
				"SearchID": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	// Defines the properties of a saved search, used when querying for saved searches
	searchQueryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Search",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: graphql.String,
			},
			"SearchKey": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	// Defines the Queries that can be made
	characterQueryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "CharacterQuery",
		Fields: graphql.Fields{
			"getCharacters": &graphql.Field{
				Type: charactersResultType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Description: "name of the character",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					name := p.Args["name"].(string)
					characters, searchID, err := svc.GetCharacters(name)
					if err != nil {
						return nil, err
					}

					return CharactersResult{Characters: characters, SearchID: searchID}, nil
				},
			},
			"getSavedSearches": &graphql.Field{
				Type: graphql.NewList(searchQueryType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					searches, err := svc.GetSavedSearches()
					if err != nil {
						return nil, err
					}
					return searches, nil
				},
			},
			"getSavedSearchesByID": &graphql.Field{
				Type: graphql.NewList(characterType),
				Args: graphql.FieldConfigArgument{
					"searchID": &graphql.ArgumentConfig{
						Description: "the search ID",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					searchID := p.Args["searchID"].(string)

					characters, err := svc.GetSavedSearchesByID(searchID)
					if err != nil {
						return nil, err
					}
					return characters, nil
				},
			},
		},
	})

	// Defines the Mutations that can be made
	characterMutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "CharacterMutation",
		Fields: graphql.Fields{
			"saveSearch": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"searchID": &graphql.ArgumentConfig{
						Description: "the search ID",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					searchID := p.Args["searchID"].(string)

					result, err := svc.SaveSearch(searchID)
					if err != nil {
						return nil, err
					}

					return result, nil
				},
			},
		},
	})

	// Combines the Queries and Mutations into a Schema
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    characterQueryType,
		Mutation: characterMutationType,
	})

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   cfg.Pretty,
		GraphiQL: cfg.GraphiQL,
	})

	return h
}
