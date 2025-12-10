package routes

import (
	"net/http"
	"strconv"
	"strings"
	"godoor/models"
	"godoor/serializers"
	"github.com/gin-gonic/gin"
)

// GetRecipes returns list of recipes with sorting and pagination
func GetRecipes(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)

	// Parse query parameters
	sortOrder := c.Query("sort_order")
	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")

	// Default values
	pageSize := 20
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
	query := models.DB.Where("space_id = ?", space.ID).Preload("CreatedBy")

	// Apply sorting
	if sortOrder != "" {
		switch sortOrder {
		case "-rating":
			// Sort by rating descending (nulls last)
			query = query.Order("rating DESC NULLS LAST")
		case "rating":
			query = query.Order("rating ASC NULLS LAST")
		case "-name":
			query = query.Order("name DESC")
		case "name":
			query = query.Order("name ASC")
		case "-created_at":
			query = query.Order("created_at DESC")
		case "created_at":
			query = query.Order("created_at ASC")
		default:
			query = query.Order("created_at DESC") // Default sort
		}
	} else {
		query = query.Order("created_at DESC") // Default sort
	}

	// Get total count
	var totalCount int64
	query.Model(&models.Recipe{}).Count(&totalCount)

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute query
	var recipes []models.Recipe
	if err := query.Find(&recipes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes"})
		return
	}

	// Serialize response using paginated serializer
	response := serializers.SerializeRecipes(recipes, int(totalCount))
	c.JSON(http.StatusOK, response)
}

// GetRecipe returns a specific recipe
func GetRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var recipe models.Recipe
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).Preload("CreatedBy").First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	response := serializers.SerializeRecipeOverview(&recipe)
	c.JSON(http.StatusOK, response)
}

// CreateRecipeRequest represents recipe creation request
type CreateRecipeRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	Servings    int     `json:"servings,omitempty"`
}

// CreateRecipe creates a new recipe
func CreateRecipe(c *gin.Context) {
	var req CreateRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user := c.MustGet("user").(*models.User)
	space := c.MustGet("space").(*models.Space)

	if req.Servings <= 0 {
		req.Servings = 1
	}

	recipe, err := models.CreateRecipe(user, space, req.Name, req.Description, req.Servings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe"})
		return
	}

	if err := models.DB.Preload("CreatedBy").First(recipe, recipe.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload recipe"})
		return
	}

	response := serializers.SerializeRecipeOverview(recipe)
	c.JSON(http.StatusCreated, response)
}

// UpdateRecipeRequest represents recipe update request
type UpdateRecipeRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Servings    *int    `json:"servings,omitempty"`
}

// UpdateRecipe updates a recipe
func UpdateRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	var req UpdateRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	currentUser := c.MustGet("user").(*models.User)

	var recipe models.Recipe
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	// Check permissions (only creator can update)
	if recipe.CreatedByID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Update fields
	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		recipe.Name = *req.Name
	}
	if req.Description != nil {
		recipe.Description = req.Description
	}
	if req.Servings != nil && *req.Servings > 0 {
		recipe.Servings = *req.Servings
	}

	if err := models.DB.Save(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipe"})
		return
	}

	if err := models.DB.Preload("CreatedBy").First(&recipe, recipe.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload recipe"})
		return
	}

	response := serializers.SerializeRecipeOverview(&recipe)
	c.JSON(http.StatusOK, response)
}

// DeleteRecipe deletes a recipe
func DeleteRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	currentUser := c.MustGet("user").(*models.User)

	var recipe models.Recipe
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	// Check permissions (only creator can delete)
	if recipe.CreatedByID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := models.DB.Delete(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
