package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"agenticify-3906125248/internal/models"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	// The original bug description mentioned a nil pointer dereference, but the provided main.go
	// initializes the response correctly. This implementation also ensures proper initialization.
	response := models.PingResponse{Message: "pong"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status:  "ok",
		Version: "1.0.0", // Example version, could be dynamic
		Uptime:  time.Now().Unix(), // Example uptime
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
