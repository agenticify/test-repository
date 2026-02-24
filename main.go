package main

import (
	"log"
	"net/http"
	"os"

	"test-repository/internal/handlers"
	"test-repository/internal/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", handlers.PingHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Apply middleware
	wrappedMux := middleware.Logger(middleware.Recovery(mux))

	log.Printf("Server is starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, wrappedMux); err != nil {
		log.Fatal(err)
	}
}
