package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestPingHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	pingHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var resp PingResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "pong" {
		t.Errorf("expected message \"pong\", got %q", resp.Message)
	}
}

func TestUpHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/up", nil)
	rec := httptest.NewRecorder()

	upHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var resp UpResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}

	if resp.Timestamp <= 0 {
		t.Errorf("expected positive timestamp, got %d", resp.Timestamp)
	}
}

func TestStatsHandler_Concurrent(t *testing.T) {
	// Reset counter
	requestCount.Store(0)

	var wg sync.WaitGroup
	n := 100

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, "/stats", nil)
			rec := httptest.NewRecorder()
			statsHandler(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
			}
		}()
	}
	wg.Wait()

	count := requestCount.Load()
	if count != int64(n) {
		t.Errorf("expected request count %d, got %d", n, count)
	}
}

func TestDeleteUserHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/delete-user?id=1", nil)
	rec := httptest.NewRecorder()

	deleteUserHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestDeleteUserHandler_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/delete-user?id=abc", nil)
	rec := httptest.NewRecorder()

	deleteUserHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
