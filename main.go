package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"agenticify-3906125248/internal/handlers"
)

// LoggerMiddleware logs incoming HTTP requests.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}

// RecoveryMiddleware recovers from panics and logs the error.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Configuration Management: Read port from environment variable, default to 3000
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || port == 0 {
		port = 3000 // Default port
	}
	addr := ":" + strconv.Itoa(port)

	// Create a custom http.ServeMux
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/ping", handlers.PingHandler)
	mux.HandleFunc("/health", handlers.HealthHandler) // Added health endpoint

	// Apply middleware
	var handler http.Handler = mux
	handler = LoggerMiddleware(handler)
	handler = RecoveryMiddleware(handler)

	log.Printf("Server is starting on port %d...", port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
