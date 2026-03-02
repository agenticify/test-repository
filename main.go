package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type PingResponse struct {
	Message string `json:"message"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := PingResponse{Message: "pong"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	healthHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rec.Code)
	}
	if rec.Body.String() != "OK" {
		t.Errorf("expected body 'OK'; got %v", rec.Body.String())
	}
}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/health", healthHandler)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
