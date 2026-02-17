package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
)

var (
	db           *sql.DB
	requestCount atomic.Int64
)

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

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	rows, err := db.Query("SELECT id, name, email FROM users WHERE name = $1", username)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
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

func statsHandler(w http.ResponseWriter, r *http.Request) {
	count := requestCount.Add(1)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"total_requests": count,
	}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		log.Printf("failed to delete user %d: %v", userID, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get rows affected: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "user deleted",
	}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// Critical: Command injection — user input passed directly to shell
func debugHandler(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, string(out))
}

// Warning: Unbounded memory — no limit on request body size
func importUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []User
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	for _, u := range users {
		_, err := db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", u.Name, u.Email)
		if err != nil {
			log.Printf("failed to insert user %s: %v", u.Name, err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"imported": len(users),
	})
}

func main() {
	initDB()

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/up", upHandler)
	http.HandleFunc("/users", getUserHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/delete-user", deleteUserHandler)
	http.HandleFunc("/debug", debugHandler)
	http.HandleFunc("/import-users", importUsersHandler)

	log.Println("Server is starting on port 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
