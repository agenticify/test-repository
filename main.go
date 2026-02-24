package main

import (
	"log"
	"net/http"
	"os"

	"/internal/handlers"
	"/internal/middleware"
)

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/ping", handlers.PingHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Apply middleware
	handledMux := middleware.Logger(middleware.Recovery(mux))

	// Get port from environment variable or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := ":" + port

	log.Printf("Server is starting on port %s...", port)
	if err := http.ListenAndServe(addr, handledMux); err != nil {
		log.Fatal(err)
	}
}
