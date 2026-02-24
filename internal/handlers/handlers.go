package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/your-username/your-project-name/internal/models"
)

// PingHandler responds to ping requests.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	response := models.PingResponse{Message: "pong"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HealthHandler responds to health check requests.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// The meta map was identified as dead code in the description.
	// We will create a response without it.
	resp := map[string]string{"status": "ok"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
