package main

import (
	"log"
	"net/http"
	"os"

	"test-repository/internal/handlers"
	"test-repository/internal/middleware"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", handlers.PingHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Apply middleware
	wrappedMux := middleware.Logger(middleware.Recovery(mux))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	addr := ":" + port
	log.Printf("Server is starting on port %s...", port)
	if err := http.ListenAndServe(addr, wrappedMux); err != nil {
		log.Fatal(err)
	}
}
