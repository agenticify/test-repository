package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/justinas/alice"
	"test-repository/internal/handlers"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	// Read port from environment variable, default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Initialize handlers
	h := handlers.NewHandlers(logger)

	// Apply middleware chain
	chain := alice.New(loggingMiddleware(logger), recoveryMiddleware(logger))

	// Register handlers with the middleware chain
	mux.Handle("/ping", chain.ThenFunc(h.PingHandler))
	mux.Handle("/health", chain.ThenFunc(h.HealthHandler))

	// Start the server
	serverAddr := fmt.Sprintf(":%s", port)
	logger.Printf("Server is starting on port %s...", port)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}

// loggingMiddleware logs details of incoming requests.
func loggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

// recoveryMiddleware recovers from panics and logs the error.
func recoveryMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					logger.Printf("Panic: %v
%s", err, debug.Stack())
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
