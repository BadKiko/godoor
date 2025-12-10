package middleware

import (
	"net/http"
	"strings"
	"godoor/models"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates access tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Find token in database
		var accessToken models.AccessToken
		if err := models.DB.Where("token = ?", token).First(&accessToken).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if token is expired
		// TODO: implement expiration check

		// Load user
		var user models.User
		if err := models.DB.First(&user, accessToken.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Get active space for user
		activeSpace := user.GetActiveSpace()
		if activeSpace == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No active space found"})
			c.Abort()
			return
		}

		// Set user and space in context
		c.Set("user", &user)
		c.Set("space", activeSpace)
		c.Next()
	}
}
