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
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName"`
	ContactNo  string    `json:"contactNo"`
	Email      string    `json:"email"`
	DOB        time.Time `json:"dob"` // Use string or time.Time for date
	Address    string    `json:"address"`
	State      string    `json:"state"`
	City       string    `json:"city"`
	Pin        string    `json:"pin"`
	Aadhar     string    `json:"aadhar"`
	Pan        string    `json:"pan"`
	Active     bool      `json:"active"` // true for active, false for inactive
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
	r.Use(CORSMiddleware())

	// Define the handler for registration
	r.POST("/api/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		query := `INSERT INTO divine_users (
            username, password, firstName, middleName, lastName, contactNo,
            email, dob, address, state, city, pin, aadhar, pan, active) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = db.Exec(
			query, user.Username, user.Password, user.FirstName, user.MiddleName, user.LastName,
			user.ContactNo, user.Email, user.DOB, user.Address, user.State, user.City,
			user.Pin, user.Aadhar, user.Pan, user.Active,
		)
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
		err := db.QueryRow("SELECT password FROM divine_users WHERE email = ?", user.Email).Scan(&storedPassword)
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
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // Handle preflight requests
			return
		}

		c.Next()
	}
}
