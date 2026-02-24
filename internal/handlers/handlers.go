package handlers

import (
	"encoding/json"
	"net/http"

	"test-repository/internal/models"
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

// HealthHandler provides a simple health check endpoint.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp := models.JSONResponse{
		Status:  http.StatusOK,
		Message: "Service is healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
