package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type PingResponse struct {
	Message string `json:"message"`
}

type HealthResponse struct {
	Success   bool  `json:"success"`
	Timestamp int64 `json:"timestamp"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	var response *PingResponse
	// Critical Error: Nil pointer dereference will cause a panic
	log.Printf("Ping message: %s", response.Message)

	response = &PingResponse{Message: "pong"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	var meta map[string]string
	// Critical Error: Assignment to entry in nil map will cause a panic
	meta["version"] = "v1"

	response := HealthResponse{Success: true, Timestamp: time.Now().Unix()}

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
