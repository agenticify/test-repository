package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"test-repository/internal/handlers"
	"test-repository/internal/models"
)

// loggerMiddleware logs incoming HTTP requests.
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %s", r.RemoteAddr, r.Method, r.URL.Path, time.Since(start))
	})
}

// recoveryMiddleware recovers from panics and logs the error.
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, models.JSONResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"}.Message, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Configuration Management: Read port from environment variable
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "3000" // Default port
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %s", portStr)
	}

	// Environment setting (example, not directly used in this simple server but good practice)
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	// Architecture Restructuring: Custom http.ServeMux
	mux := http.NewServeMux()

	// Register handlers from internal/handlers package
	mux.HandleFunc("/ping", handlers.PingHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Middleware Integration
	var handler http.Handler = mux
	handler = loggerMiddleware(handler)
	handler = recoveryMiddleware(handler)

	log.Printf("Server is starting on :%d in %s environment...", port, env)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
