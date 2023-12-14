package internal

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/rs/cors"
)

type Server struct {
	Schema *graphql.Schema
	Server *http.Server
}

type ServerConfig struct {
	Port string
}

// NewServer returns a new HTTP server
func NewServer(cfg ServerConfig, h *handler.Handler) *http.Server {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowCredentials: true,
		AllowedMethods:   []string{"OPTIONS", "POST"},
	})

	// Create a new HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: c.Handler(h),
	}

	return srv
}
