package routes

import (
	"net/http"
	"strconv"
	"godoor/models"
	"godoor/serializers"
	"github.com/gin-gonic/gin"
)

// GetCookLogs returns list of cook logs with pagination and recipe filtering
func GetCookLogs(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)

	// Parse query parameters
	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")
	recipeIDStr := c.Query("recipe")

	// Default values
	pageSize := 50
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Build query
	query := models.DB.Where("space_id = ?", space.ID).Preload("Recipe").Preload("CreatedBy")

	// Filter by recipe if specified
	if recipeIDStr != "" {
		if recipeID, err := strconv.Atoi(recipeIDStr); err == nil {
			query = query.Where("recipe_id = ?", uint(recipeID))
		}
	}

	// Get total count
	var totalCount int64
	query.Model(&models.CookLog{}).Count(&totalCount)

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute query
	var cookLogs []models.CookLog
	if err := query.Find(&cookLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cook logs"})
		return
	}

	// Serialize response
	response := serializers.SerializeCookLogs(cookLogs, int(totalCount))
	c.JSON(http.StatusOK, response)
}

// GetCookLog returns a specific cook log
func GetCookLog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cook log ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var cookLog models.CookLog
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).Preload("Recipe").Preload("CreatedBy").First(&cookLog).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cook log not found"})
		return
	}

	response := serializers.SerializeCookLog(&cookLog)
	c.JSON(http.StatusOK, response)
}

// CreateCookLogRequest represents cook log creation request
type CreateCookLogRequest struct {
	RecipeID uint    `json:"recipe" binding:"required"`
	Servings *int    `json:"servings,omitempty"`
	Rating   *int    `json:"rating,omitempty"`
	Comment  *string `json:"comment,omitempty"`
}

// CreateCookLog creates a new cook log
func CreateCookLog(c *gin.Context) {
	var req CreateCookLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user := c.MustGet("user").(*models.User)
	space := c.MustGet("space").(*models.Space)

	// Verify recipe exists and belongs to space
	var recipe models.Recipe
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, req.RecipeID).First(&recipe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Recipe not found"})
		return
	}

	cookLog := models.CookLog{
		RecipeID:   req.RecipeID,
		Recipe:     recipe,
		Servings:   req.Servings,
		Rating:     req.Rating,
		Comment:    req.Comment,
		CreatedByID: user.ID,
		CreatedBy:   *user,
		SpaceID:     space.ID,
		Space:       *space,
	}

	if err := models.DB.Create(&cookLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cook log"})
		return
	}

	if err := models.DB.Preload("Recipe").Preload("CreatedBy").First(&cookLog, cookLog.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload cook log"})
		return
	}

	response := serializers.SerializeCookLog(&cookLog)
	c.JSON(http.StatusCreated, response)
}

// UpdateCookLogRequest represents cook log update request
type UpdateCookLogRequest struct {
	Servings *int    `json:"servings,omitempty"`
	Rating   *int    `json:"rating,omitempty"`
	Comment  *string `json:"comment,omitempty"`
}

// UpdateCookLog updates a cook log
func UpdateCookLog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cook log ID"})
		return
	}

	var req UpdateCookLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	currentUser := c.MustGet("user").(*models.User)

	var cookLog models.CookLog
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&cookLog).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cook log not found"})
		return
	}

	// Check permissions (only creator can update)
	if cookLog.CreatedByID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Update fields
	if req.Servings != nil {
		cookLog.Servings = req.Servings
	}
	if req.Rating != nil {
		cookLog.Rating = req.Rating
	}
	if req.Comment != nil {
		cookLog.Comment = req.Comment
	}

	if err := models.DB.Save(&cookLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cook log"})
		return
	}

	if err := models.DB.Preload("Recipe").Preload("CreatedBy").First(&cookLog, cookLog.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload cook log"})
		return
	}

	response := serializers.SerializeCookLog(&cookLog)
	c.JSON(http.StatusOK, response)
}

// DeleteCookLog deletes a cook log
func DeleteCookLog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cook log ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	currentUser := c.MustGet("user").(*models.User)

	var cookLog models.CookLog
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&cookLog).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cook log not found"})
		return
	}

	// Check permissions (only creator can delete)
	if cookLog.CreatedByID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := models.DB.Delete(&cookLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cook log"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
