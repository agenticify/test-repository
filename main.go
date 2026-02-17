package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var startTime = time.Now()

type PingResponse struct {
	Message string `json:"message"`
}

type HealthResponse struct {
	Status    string  `json:"status"`
	Uptime    float64 `json:"uptime"`
	Version   string  `json:"version"`
	Timestamp string  `json:"timestamp"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := PingResponse{Message: "pong"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime).Seconds()
	response := HealthResponse{
		Status:    "ok",
		Uptime:    uptime,
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/health", healthHandler)

	log.Println("Server is starting on port 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
