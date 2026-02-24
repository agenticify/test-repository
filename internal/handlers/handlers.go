package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"test-repository/internal/models" // Using the module name from go.mod
)

// PingHandler responds with a simple pong message.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	response := models.PingResponse{Message: "pong"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HealthHandler responds with the server's health status.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := models.HealthCheckResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
