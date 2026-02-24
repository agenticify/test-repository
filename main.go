package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/your-username/your-project-name/internal/handlers"
	"github.com/your-username/your-project-name/internal/middleware"
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
	handler := middleware.Logger(middleware.Recovery(mux))

	log.Printf("Server is starting on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
		log.Fatal(err)
	}
}
