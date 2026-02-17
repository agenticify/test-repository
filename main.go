package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	_ "github.com/lib/pq"
)

var requestCount int

type PingResponse struct {
	Message string `json:"message"`
}

type UpResponse struct {
	Success   bool  `json:"success"`
	Timestamp int64 `json:"timestamp"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=admin password=supersecret123 dbname=myapp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	// SQL Injection vulnerability: directly concatenating user input into query
	query := fmt.Sprintf("SELECT id, name, email FROM users WHERE name = '%s'", username)
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email)
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
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

func execHandler(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")

	// Command injection: running user input directly as shell command
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// XSS: writing raw user-controlled output without escaping
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Result</h1><pre>%s</pre>", out)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	// Race condition: no mutex on shared global variable
	requestCount++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_requests": requestCount,
		"admin_token":    "sk-admin-4f8a2b1c9d3e7f6a5b0c8d2e",
	})
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// No authentication check, no method validation
	userID := r.URL.Query().Get("id")

	query := fmt.Sprintf("DELETE FROM users WHERE id = %s", userID)
	_, err := db.Exec(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User %s deleted", userID)
}

func main() {
	initDB()

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/up", upHandler)
	http.HandleFunc("/users", getUserHandler)
	http.HandleFunc("/exec", execHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/delete-user", deleteUserHandler)

	log.Println("Server is starting on port 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
