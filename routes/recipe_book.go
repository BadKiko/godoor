package routes

import (
	"net/http"
	"strconv"
	"godoor/models"
	"godoor/serializers"
	"github.com/gin-gonic/gin"
)

// GetRecipeBooks returns list of recipe books for current user
func GetRecipeBooks(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)

	var books []models.RecipeBook
	if err := models.DB.Where("space_id = ?", space.ID).Preload("CreatedBy").Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipe books"})
		return
	}

	response := serializers.SerializeRecipeBooks(books)
	c.JSON(http.StatusOK, response)
}

// GetRecipeBook returns a specific recipe book
func GetRecipeBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe book ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var book models.RecipeBook
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).Preload("CreatedBy").First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe book not found"})
		return
	}

	response := serializers.SerializeRecipeBook(&book)
	c.JSON(http.StatusOK, response)
}

// CreateRecipeBookRequest represents recipe book creation request
type CreateRecipeBookRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

// CreateRecipeBook creates a new recipe book
func CreateRecipeBook(c *gin.Context) {
	var req CreateRecipeBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user := c.MustGet("user").(*models.User)
	space := c.MustGet("space").(*models.Space)

	book, err := models.CreateRecipeBook(user, space, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe book"})
		return
	}

	response := serializers.SerializeRecipeBook(book)
	c.JSON(http.StatusCreated, response)
}

// UpdateRecipeBookRequest represents recipe book update request
type UpdateRecipeBookRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Order       *int    `json:"order,omitempty"`
}

// UpdateRecipeBook updates a recipe book
func UpdateRecipeBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe book ID"})
		return
	}

	var req UpdateRecipeBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	currentUser := c.MustGet("user").(*models.User)

	var book models.RecipeBook
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe book not found"})
		return
	}

	// Check permissions (only creator can update)
	if book.CreatedByID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Update fields
	if req.Name != nil {
		book.Name = *req.Name
	}
	if req.Description != nil {
		book.Description = *req.Description
	}
	if req.Order != nil {
		book.Order = *req.Order
	}

	if err := models.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipe book"})
		return
	}

	if err := models.DB.Preload("CreatedBy").First(&book, book.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload recipe book"})
		return
	}

	response := serializers.SerializeRecipeBook(&book)
	c.JSON(http.StatusOK, response)
}

// DeleteRecipeBook deletes a recipe book
func DeleteRecipeBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe book ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	currentUser := c.MustGet("user").(*models.User)

	var book models.RecipeBook
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe book not found"})
		return
	}

	// Check permissions (only creator can delete)
	if book.CreatedByID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := models.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe book"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
