package internal

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Server struct {
	Schema *graphql.Schema
	Server *http.Server
}

type ServerConfig struct {
	Pretty   bool
	GraphiQL bool
	Port     string
}

var (
	characterType *graphql.Object
)

// initializes characterType when the package is imported
func init() {
	characterType = graphql.NewObject(graphql.ObjectConfig{
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
						return character.CarModels, nil
					}
					return []interface{}{}, nil
				},
			},
			"isBookmarked": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Is the character bookmarked.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if character, ok := p.Source.(Character); ok {
						return character.IsBookmarked, nil
					}
					return false, nil
				},
			},
		},
	})
}

func NewServer(cfg ServerConfig, svc CharacterService) Server {
	characterQueryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "CharacterQuery",
		Fields: graphql.Fields{
			"characters": &graphql.Field{
				Type: graphql.NewList(characterType),
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Description: "name of the character",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					name := p.Args["name"].(string)
					characters, err := svc.GetCharacters(name)
					if err != nil {
						return nil, err
					}

					return characters, nil
				},
			},
			"bookmarks": &graphql.Field{
				Type: graphql.NewList(characterType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					characters, err := svc.GetSavedSearches()
					if err != nil {
						return nil, err
					}
					return characters, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: characterQueryType,
	})

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   cfg.Pretty,
		GraphiQL: cfg.GraphiQL,
	})

	// Create a new HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: h,
	}

	return Server{
		Schema: &schema,
		Server: srv,
	}
}
