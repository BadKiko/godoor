package routes

import (
	"net/http"
	"strconv"
	"godoor/models"
	"godoor/serializers"
	"github.com/gin-gonic/gin"
)

// GetSpaces returns list of spaces for current user
func GetSpaces(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var userSpaces []models.UserSpace
	if err := models.DB.Where("user_id = ?", user.ID).Preload("Space").Find(&userSpaces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch spaces"})
		return
	}

	spaces := make([]models.Space, len(userSpaces))
	for i, us := range userSpaces {
		spaces[i] = us.Space
	}

	c.JSON(http.StatusOK, serializers.SerializeSpaces(spaces))
}

// CreateSpaceRequest represents space creation request
type CreateSpaceRequest struct {
	Name *string `json:"name,omitempty"`
}

// CreateSpace creates a new space for user
func CreateSpace(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var req CreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Check space limit (simplified - no limit for now)
	// TODO: implement space limits like in Tandoor

	userSpace, err := models.CreateSpaceForUser(user, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create space"})
		return
	}

	response := serializers.SerializeSpace(&userSpace.Space)
	createdBy := serializers.SerializeUser(user)
	response.CreatedBy = &createdBy

	c.JSON(http.StatusCreated, response)
}

// GetSpace returns a specific space
func GetSpace(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid space ID"})
		return
	}

	user := c.MustGet("user").(*models.User)

	var userSpace models.UserSpace
	if err := models.DB.Where("user_id = ? AND space_id = ?", user.ID, uint(id)).Preload("Space.CreatedBy").First(&userSpace).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Space not found"})
		return
	}

	response := serializers.SerializeSpace(&userSpace.Space)
	createdBy := serializers.SerializeUser(&userSpace.Space.CreatedBy)
	response.CreatedBy = &createdBy

	c.JSON(http.StatusOK, response)
}

// UpdateSpaceRequest represents space update request
type UpdateSpaceRequest struct {
	Name    *string `json:"name,omitempty"`
	Message *string `json:"message,omitempty"`
}

// UpdateSpace updates a space
func UpdateSpace(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid space ID"})
		return
	}

	var req UpdateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user := c.MustGet("user").(*models.User)

	var space models.Space
	if err := models.DB.First(&space, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Space not found"})
		return
	}

	// Check ownership
	if space.CreatedByID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only space owner can update space"})
		return
	}

	// Update fields
	if req.Name != nil {
		space.Name = *req.Name
	}
	if req.Message != nil {
		space.Message = *req.Message
	}

	if err := models.DB.Save(&space).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update space"})
		return
	}

	response := serializers.SerializeSpace(&space)
	createdBySerializer := serializers.SerializeUser(user)
	response.CreatedBy = &createdBySerializer

	c.JSON(http.StatusOK, response)
}

// GetCurrentSpace returns the current active space
func GetCurrentSpace(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)

	response := serializers.SerializeSpace(space)
	// TODO: add created_by if needed

	c.JSON(http.StatusOK, response)
}
