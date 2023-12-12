package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"alvinlucillo/swapi-app/internal"
)

func main() {

	var (
		prettyFlag   = flag.Bool("pretty", true, "Enable/disable pretty formatting")
		graphiQLFlag = flag.Bool("graphiql", true, "Enable/disable GraphiQL")
		portFlag     = flag.String("port", "8080", "Port number for the server")
	)

	flag.Parse()

	swapiApiClient := internal.NewSWAPIClient(&http.Client{}, "https://swapi.dev/api")

	svc := internal.NewService(swapiApiClient)

	srv := internal.NewServer(internal.ServerConfig{
		GraphiQL: *graphiQLFlag,
		Port:     *portFlag,
		Pretty:   *prettyFlag,
	}, svc)

	// Start the server in a separate goroutine
	go func() {
		fmt.Println("Server is running on http://localhost:8080/graphql")
		if err := srv.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
	if err := srv.Server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}
