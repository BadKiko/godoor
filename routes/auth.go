package routes

import (
	"net/http"
	"godoor/models"
	"github.com/gin-gonic/gin"
)

// LoginRequest represents login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthTokenResponse represents the response from token authentication (matching Tandoor)
type AuthTokenResponse struct {
	ID       uint   `json:"id"`
	Token    string `json:"token"`
	Scope    string `json:"scope"`
	Expires  string `json:"expires"` // TODO: format as ISO string
	UserID   uint   `json:"user_id"`
	Test     uint   `json:"test"`     // This field exists in Tandoor response
}

// AuthenticateUser handles user authentication
func AuthenticateUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Find user by username
	var user models.User
	if err := models.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create or get existing token
	token, err := models.CreateTokenForUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	// Return response matching Tandoor format
	response := AuthTokenResponse{
		ID:      token.ID,
		Token:   token.Token,
		Scope:   token.Scope,
		Expires: token.Expires.Format("2006-01-02T15:04:05Z07:00"),
		UserID:  token.UserID,
		Test:    token.UserID, // This field is in Tandoor response
	}

	c.JSON(http.StatusOK, response)
}
