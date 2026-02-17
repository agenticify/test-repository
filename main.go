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

type UpResponse struct {
	Success   bool  `json:"success"`
	Timestamp int64 `json:"timestamp"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := PingResponse{Message: "pong"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func upHandler(w http.ResponseWriter, r *http.Request) {
	response := UpResponse{
		Success:   true,
		Timestamp: time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/up", upHandler)

	log.Println("Server is starting on port 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
