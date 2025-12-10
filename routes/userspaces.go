package routes

import (
	"net/http"
	"strconv"
	"godoor/models"
	"godoor/serializers"
	"github.com/gin-gonic/gin"
)

// GetUserSpaces returns user-space relationships
func GetUserSpaces(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)

	var userSpaces []models.UserSpace
	if err := models.DB.Where("space_id = ?", space.ID).Preload("User").Preload("Space").Find(&userSpaces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user spaces"})
		return
	}

	c.JSON(http.StatusOK, serializers.SerializeUserSpaces(userSpaces))
}

// GetUserSpace returns a specific user-space relationship
func GetUserSpace(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user space ID"})
		return
	}

	var userSpace models.UserSpace
	if err := models.DB.Preload("User").Preload("Space").First(&userSpace, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User space not found"})
		return
	}

	c.JSON(http.StatusOK, serializers.SerializeUserSpace(&userSpace))
}

// UpdateUserSpaceRequest represents user space update request
type UpdateUserSpaceRequest struct {
	Active       *bool   `json:"active,omitempty"`
	InternalNote *string `json:"internal_note,omitempty"`
}

// UpdateUserSpace updates a user-space relationship
func UpdateUserSpace(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user space ID"})
		return
	}

	var req UpdateUserSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	currentUser := c.MustGet("user").(*models.User)

	var userSpace models.UserSpace
	if err := models.DB.Preload("User").Preload("Space").First(&userSpace, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User space not found"})
		return
	}

	// Check permissions (only space owner can change other users' settings)
	if userSpace.User.ID != currentUser.ID && (userSpace.Space.CreatedByID == nil || *userSpace.Space.CreatedByID != currentUser.ID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Update fields
	if req.Active != nil {
		userSpace.Active = *req.Active
	}
	if req.InternalNote != nil {
		userSpace.InternalNote = req.InternalNote
	}

	if err := models.DB.Save(&userSpace).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user space"})
		return
	}

	c.JSON(http.StatusOK, serializers.SerializeUserSpace(&userSpace))
}
