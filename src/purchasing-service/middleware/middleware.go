package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"purchasing_service/structs"
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
	var u structs.Sanco_Users

	// Bind the incoming JSON request to the user struct
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch user by username
	var user structs.Sanco_Users
	if err := DB.Where("username = ?", u.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Compare the provided password with the stored hashed password
	if CheckPasswordHash(u.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create JWT token upon successful login
	tokenString, err := CreateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	// Return the generated token in the response
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
