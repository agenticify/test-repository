package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"refactor/internal/models" // Assuming "refactor" is the module name
)

// Handlers struct holds dependencies for handlers.
type Handlers struct {
	Logger *log.Logger
}

// NewHandlers creates and returns a new Handlers instance.
func NewHandlers(logger *log.Logger) *Handlers {
	return &Handlers{Logger: logger}
}

// PingHandler responds with a simple pong message.
func (h *Handlers) PingHandler(w http.ResponseWriter, r *http.Request) {
	resp := models.PingResponse{Message: "pong"}
	h.sendJSONResponse(w, http.StatusOK, resp)
}

// HealthHandler responds with a health status.
func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp := models.Response{
		Status:  http.StatusOK,
		Message: "Service is healthy",
	}
	h.sendJSONResponse(w, http.StatusOK, resp)
}

// sendJSONResponse is a helper to send JSON responses.
func (h *Handlers) sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.Logger.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
