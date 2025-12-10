package routes

import (
	"net/http"
	"strconv"
	"godoor/models"
	"godoor/serializers"
	"github.com/gin-gonic/gin"
)

// GetKeywords returns list of keywords for current user
func GetKeywords(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)

	var keywords []models.Keyword
	if err := models.DB.Where("space_id = ?", space.ID).Find(&keywords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch keywords"})
		return
	}

	response := serializers.SerializeKeywords(keywords)
	c.JSON(http.StatusOK, response)
}

// GetKeyword returns a specific keyword
func GetKeyword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid keyword ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var keyword models.Keyword
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&keyword).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Keyword not found"})
		return
	}

	response := serializers.SerializeKeyword(&keyword)
	c.JSON(http.StatusOK, response)
}

// CreateKeywordRequest represents keyword creation request
type CreateKeywordRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

// CreateKeyword creates a new keyword
func CreateKeyword(c *gin.Context) {
	var req CreateKeywordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	keyword, err := models.CreateKeyword(space, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create keyword"})
		return
	}

	response := serializers.SerializeKeyword(keyword)
	c.JSON(http.StatusCreated, response)
}

// UpdateKeywordRequest represents keyword update request
type UpdateKeywordRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdateKeyword updates a keyword
func UpdateKeyword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid keyword ID"})
		return
	}

	var req UpdateKeywordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var keyword models.Keyword
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&keyword).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Keyword not found"})
		return
	}

	// Update fields
	if req.Name != nil {
		keyword.Name = *req.Name
	}
	if req.Description != nil {
		keyword.Description = *req.Description
	}

	if err := models.DB.Save(&keyword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update keyword"})
		return
	}

	response := serializers.SerializeKeyword(&keyword)
	c.JSON(http.StatusOK, response)
}

// DeleteKeyword deletes a keyword
func DeleteKeyword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid keyword ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var keyword models.Keyword
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&keyword).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Keyword not found"})
		return
	}

	if err := models.DB.Delete(&keyword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete keyword"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
