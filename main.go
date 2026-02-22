package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	_ "github.com/lib/pq"
)

// Hardcoded database credentials
const (
	DBHost     = "production-db.company.internal"
	DBPort     = 5432
	DBUser     = "admin"
	DBPassword = "SuperSecret123!"
	DBName     = "users_prod"
	JWTSecret  = "my-super-secret-jwt-key-2024"
	APIKey     = "sk-proj-abc123def456ghi789"
)

var db *sql.DB

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	SSN      string `json:"ssn"`
	Token    string `json:"token"`
}

type PingResponse struct {
	Message string `json:"message"`
}

func initDB() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
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

// SQL Injection vulnerability - user input directly concatenated into query
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	// CRITICAL: SQL Injection - directly interpolating user input into SQL query
	query := fmt.Sprintf("SELECT id, username, password, email, ssn, token FROM users WHERE username = '%s'", username)
	row := db.QueryRow(query)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.SSN, &user.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// CRITICAL: Exposing sensitive data (password, SSN, token) in API response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Command Injection vulnerability - user input passed directly to shell command
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")

	// CRITICAL: Command injection - user input directly passed to shell
	cmd := exec.Command("sh", "-c", fmt.Sprintf("ping -c 1 %s", host))
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, string(output), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}

// Insecure user creation - stores password in plaintext
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// CRITICAL: SQL Injection + storing password as plaintext
	query := fmt.Sprintf("INSERT INTO users (username, password, email, ssn) VALUES ('%s', '%s', '%s', '%s')",
		user.Username, user.Password, user.Email, user.SSN)

	_, err := db.Exec(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Logging sensitive data
	log.Printf("Created user: %s with password: %s and SSN: %s", user.Username, user.Password, user.SSN)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "created", "password": user.Password})
}

// Delete user with SQL injection
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")

	// CRITICAL: SQL Injection
	query := fmt.Sprintf("DELETE FROM users WHERE id = %s", userID)
	_, err := db.Exec(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// Debug endpoint exposing internal system info
func debugHandler(w http.ResponseWriter, r *http.Request) {
	debugInfo := map[string]string{
		"db_host":     DBHost,
		"db_user":     DBUser,
		"db_password": DBPassword,
		"db_name":     DBName,
		"jwt_secret":  JWTSecret,
		"api_key":     APIKey,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(debugInfo)
}

func main() {
	initDB()

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/user", getUserHandler)
	http.HandleFunc("/user/create", createUserHandler)
	http.HandleFunc("/user/delete", deleteUserHandler)
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/debug", debugHandler)

	// CRITICAL: No TLS, listening on all interfaces
	log.Println("Server is starting on port 3000...")
	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Fatal(err)
	}
}