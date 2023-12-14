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
	Port string
}

// NewServer returns a new HTTP server
func NewServer(cfg ServerConfig, h *handler.Handler) *http.Server {

	// Create a new HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: h,
	}

	return srv
}
