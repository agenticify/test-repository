package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"test-repository/internal/handlers"
	"test-repository/internal/models"
)

// Logger middleware logs incoming requests.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Recovery middleware recovers from panics and logs the error, returning a 500.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				handlers.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Configuration Management: Enable port and environment settings to be read via environment variables.
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port
	}

	mux := http.NewServeMux()

	// Register handlers from the internal/handlers package
	mux.HandleFunc("/ping", handlers.PingHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Set a custom NotFoundHandler for the mux
	mux.HandleFunc("/", handlers.NotFoundHandler)

	// Apply middleware
	var handler http.Handler = mux
	handler = Logger(handler)
	handler = Recovery(handler)

	log.Printf("Server is starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
