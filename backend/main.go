package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

type User struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	EmpName    string `json:"empName,omitempty"`
	EmpAge     int    `json:"empAge,omitempty"`
	Department string `json:"Department,omitempty"`
}

var jwtSecret = []byte("your_secret_key") // Replace with a secure key

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Connect to the database
	dbPassword := os.Getenv("DB_PASSWORD")
	db, err := sql.Open("mysql", "root:"+dbPassword+"@tcp(127.0.0.1:3306)/aayushdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// CORS setup
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})

	// Registration handler
	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		query := "INSERT INTO divine_users (username, password, empName, empAge, Department) VALUES (?, ?, ?, ?, ?)"
		_, err = db.Exec(query, user.Username, user.Password, user.EmpName, user.EmpAge, user.Department)
		if err != nil {
			http.Error(w, "Error inserting data into database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered successfully"))
	})

	// Login handler
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Check user credentials
		var storedPassword string
		err = db.QueryRow("SELECT password FROM divine_users WHERE username = ?", user.Username).Scan(&storedPassword)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		if storedPassword != user.Password {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		})
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	})

	// Wrap handler with CORS
	corsHandler := handlers.CORS(originsOk, headersOk, methodsOk)(http.DefaultServeMux)

	// Start server
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("Server running on port " + port)
	err = http.ListenAndServe(":"+port, corsHandler)
	if err != nil {
		log.Fatal(err)
	}
}
