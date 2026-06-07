package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTSecret is set from config during app initialization
var JWTSecret string

// AuthMiddleware validates JWT tokens on protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Use: Bearer <token>"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		userIDFloat, _ := claims["user_id"].(float64)
		username, _ := claims["username"].(string)
		roleIDFloat, _ := claims["role_id"].(float64)

		role := "user"
		if int(roleIDFloat) == 1 {
			role = "admin"
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized - Admin access required"})
			c.Abort()
			return
		}

		// Set user info in context for downstream handlers
		c.Set("userID", int(userIDFloat))
		c.Set("username", username)
		c.Set("role", role)

		c.Next()
	}
}
