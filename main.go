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
	addr := ":" + port

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/ping", handlers.PingHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Apply middleware
	var handler http.Handler = mux
	handler = middleware.LoggerMiddleware(handler)
	handler = middleware.RecoveryMiddleware(handler)

	log.Printf("Server is starting on port %s...", port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
