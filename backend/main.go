package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

type User struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	EmpName    string `json:"empName"`
	EmpAge     int    `json:"empAge"`
	Department string `json:"Department"`
}

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
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"}) // Allow your frontend origin

	// Define the handler for registration
	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Decode the JSON request body
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Insert the new user into the database
		query := "INSERT INTO divine_users (username, password, empName, empAge, Department) VALUES (?, ?, ?, ?, ?)"
		_, err = db.Exec(query, user.Username, user.Password, user.EmpName, user.EmpAge, user.Department)
		if err != nil {
			http.Error(w, "Error inserting data into database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered successfully"))
	})

	// Wrap your handler with CORS
	corsHandler := handlers.CORS(originsOk, headersOk, methodsOk)(http.DefaultServeMux)

	// Start server
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8000" // Default port if not specified
	}
	fmt.Println("Server running on port " + port)
	err = http.ListenAndServe(":"+port, corsHandler)
	if err != nil {
		log.Fatal(err)
	}
}
