package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	secretKey = []byte("secret-key")
)

func SetDatabase(db *gorm.DB) {
	DB = db
}

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func LoginHandler(c *gin.Context) {
	var req struct {
		SecretKey string `json:"secret_key" binding:"required"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Load expected secret from environment
	expectedSecretKey := os.Getenv("SECRET_KEY")
	if expectedSecretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server misconfiguration: SECRET_KEY not set"})
		return
	}

	// Compare provided secret key
	if req.SecretKey != expectedSecretKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid secret key"})
		return
	}

	// Generate token (can be static user or claim-based)
	tokenString, err := CreateToken("admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 403, "error": "Forbidden Missing authorization header"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("username", claims["username"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

func ProtectedHandler(c *gin.Context) {
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Welcome %s with role %s", username, role),
	})
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this resource"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares the hashed password with the plain password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
