package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	_"os"
	"time"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"cvwo-project/database"
)

type User struct {
	UserID   int	`json:"user_id"`
	Username string `json:"username"`
	Password_hash string	`json:"password_hash"`
	
}

var secretKey = []byte("my-secret-key")

type Claims struct {
	UserID   int    `json:"user_id"`
	jwt.StandardClaims
}

func Register(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password_hash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Store the user in the database
	_, err = database.DB.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", user.Username, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Generate  token
    token, err := generateToken(user.UserID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "Token": token})
}

func Login(c *gin.Context) {
	var inputUser User

	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the user from the database
	var storedUser User
	row := database.DB.QueryRow("SELECT * FROM users WHERE username = $1", inputUser.Username)
	err := row.Scan(&storedUser.UserID, &storedUser.Username, &storedUser.Password_hash)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Compare the hashed password from the database with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password_hash), []byte(inputUser.Password_hash))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate token
	token, err := generateToken(storedUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func generateToken(user_id int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token from Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if strings.HasPrefix(tokenString, bearerPrefix) {
			tokenString = strings.TrimPrefix(tokenString, bearerPrefix)
		}


		// Parse the token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		
		fmt.Println("Received token %s", tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Verify the token
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// Pass the claims to the next handler
			c.Set("claims", claims)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}
	}
}


//GetUserProfile 
func GetUserProfile(c *gin.Context) {
	// Extract claims from the context using AuthMiddleware
	claims, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Use the claims as needed (e.g., retrieve user ID)
	userID := claims.(*Claims).UserID

	// Retrieve the user profile from the database
	var userProfile User
	err:= database.DB.QueryRow("SELECT * FROM users WHERE user_id = $1", userID).Scan(
		&userProfile.UserID, &userProfile.Username, &userProfile.Password_hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	c.JSON(http.StatusOK, userProfile)
}

//Delete User (If user decides to their own account)
func DeleteUser(c *gin.Context) {
	// Extract claims from the context using AuthMiddleware
	claims, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Use the claims to get the user ID
	userID := claims.(*Claims).UserID

	// Delete the user from the database
	result, err := database.DB.Exec("DELETE FROM users WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check rows affected"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

