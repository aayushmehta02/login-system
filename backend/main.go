package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type User struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	EmpName    string `json:"empName"`
	EmpAge     int    `json:"empAge"`
	Department string `json:"Department"`
}

var jwtKey = []byte("your_secret_key")

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

	// Initialize Gin router
	r := gin.Default()

	// Define the handler for registration
	r.POST("/api/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		query := "INSERT INTO divine_users (username, password, empName, empAge, Department) VALUES (?, ?, ?, ?, ?)"
		_, err = db.Exec(query, user.Username, user.Password, user.EmpName, user.EmpAge, user.Department)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting data into database"})
			return
		}

		// Generate JWT token
		tokenString, err := generateJWT(user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "token": tokenString})
	})

	// Login handler
	r.POST("/api/login", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		var storedPassword string
		err := db.QueryRow("SELECT password FROM divine_users WHERE username = ?", user.Username).Scan(&storedPassword)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		if storedPassword != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		})
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	})

	// Start server
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8000" // Default port if not specified
	}
	fmt.Println("Server running on port " + port)
	err = r.Run(":" + port) // Start server
	if err != nil {
		log.Fatal(err)
	}
}

// Function to generate JWT
func generateJWT(username string) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), // Token expires after 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
