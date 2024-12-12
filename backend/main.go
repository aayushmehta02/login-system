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
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	Username   string    `json:"username" binding:"required"`
	Password   string    `json:"password"`
	FirstName  string    `json:"firstName" binding:"required"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName" binding:"required"`
	ContactNo  string    `json:"contactNo" binding:"required"`
	Email      string    `json:"email" binding:"required,email"`
	DOB        time.Time `json:"dob"`
	Address    string    `json:"address" binding:"required"`
	State      string    `json:"state" binding:"required"`
	City       string    `json:"city" binding:"required"`
	Pin        string    `json:"pin" binding:"required"`
	Aadhar     string    `json:"aadhar"`
	Pan        string    `json:"pan"`
	Active     bool      `json:"active"` // true for active, false for inactive
}

var jwtKey = []byte("your_secret_key")

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to the database
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD not set in environment variables")
	}

	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/aayushdb?parseTime=true", dbPassword)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()
	r.Use(CORSMiddleware())

	// Registration handler
	r.POST("/api/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Hash the password
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		// Insert user into the database
		query := `INSERT INTO divine_users (
			username, password, firstName, middleName, lastName, contactNo,
			email, dob, address, state, city, pin, aadhar, pan, active) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = db.Exec(
			query, user.Username, hashedPassword, user.FirstName, user.MiddleName, user.LastName,
			user.ContactNo, user.Email, user.DOB, user.Address, user.State, user.City,
			user.Pin, user.Aadhar, user.Pan, user.Active,
		)
		if err != nil {
			log.Printf("Error inserting data into database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting data into database"})
			return
		}

		// Generate JWT token
		tokenString, err := generateJWT(user.Username)
		if err != nil {
			log.Printf("Error generating JWT: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "token": tokenString})
	})

	// Login handler
	r.POST("/api/login", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		var storedPassword string
		var username string
		err := db.QueryRow("SELECT password, username FROM divine_users WHERE email = ?", user.Email).Scan(&storedPassword, &username)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("User not found with email: %s", user.Email)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				return
			}
			log.Printf("Error querying database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Compare hashed password with the provided password
		match, err := VerifyPassword(storedPassword, user.Password)
		if err != nil {
			log.Printf("Error verifying password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying password"})
			return
		}
		if !match {
			log.Printf("Invalid password for email: %s", user.Email)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		// Generate JWT token
		tokenString, err := generateJWT(username)
		if err != nil {
			log.Printf("Error generating JWT: %v", err)
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
		IssuedAt:  time.Now().Unix(),
		Issuer:    "your_app_name", // Optional: specify issuer
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // You can adjust the cost
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword compares a bcrypt hashed password with its possible plaintext equivalent
func VerifyPassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, err
}

// CORSMiddleware sets up CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins; consider restricting in production
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
