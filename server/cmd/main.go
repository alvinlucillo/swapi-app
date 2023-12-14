package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"alvinlucillo/swapi-app/internal"
	"alvinlucillo/swapi-app/internal/services"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Pretty   bool   `env:"PRETTY" envDefault:"true"`
	GraphiQL bool   `env:"GRAPHIQL" envDefault:"true"`
	Port     string `env:"PORT" envDefault:"8080"`
}

// Entry point of the application
func main() {
	// Load environment variables into config
	cfg, err := NewConfig()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	// Create a new SWAPI client
	swapiApiClient := services.NewSWAPIClient(&http.Client{}, "https://swapi.dev/api")
	// Create a new service
	svc, err := services.NewService(swapiApiClient)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	// Create a new handler
	h := services.NewHandler(services.HandlerConfig{Pretty: cfg.Pretty, GraphiQL: cfg.GraphiQL}, svc)
	srv := internal.NewServer(internal.ServerConfig{
		Port: cfg.Port,
	}, h)

	// Start the server in a separate goroutine
	go func() {
		fmt.Println("Server is running on http://localhost:8080/graphql")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
