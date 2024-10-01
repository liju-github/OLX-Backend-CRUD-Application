package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/liju-github/internal/auth"
	"github.com/liju-github/internal/repository"
)

func AuthMiddleware(repo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Log the authHeader to see what we're getting
		fmt.Println("Authorization Header:", authHeader)

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Extract the token after "Bearer " and ensure there are no leading/trailing spaces
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		
		// Log the token string after trimming to verify
		fmt.Println("Extracted Token:", tokenString)

		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return auth.JwtKey, nil
		})

		if err != nil || !token.Valid {
			fmt.Println("Invalid token:", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if the user exists in the database
		userEmail := claims.UserEmail
		user, err := repo.GetUserByEmail(userEmail)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User does not exist"})
			c.Abort()
			return
		}

		// Store the user email in the context
		c.Set("useremail", userEmail)
		c.Next() // Proceed to the next handler
	}
}
