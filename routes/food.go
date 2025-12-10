package routes

import (
	"net/http"
	"strconv"
	"godoor/models"
	"godoor/serializers"
	"github.com/gin-gonic/gin"
)

// GetFoods returns list of foods for current user
func GetFoods(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)

	var foods []models.Food
	query := models.DB.Where("space_id = ?", space.ID)

	// Handle search query
	if searchQuery := c.Query("query"); searchQuery != "" {
		query = query.Where("name ILIKE ?", "%"+searchQuery+"%")
	}

	var totalCount int64
	query.Model(&models.Food{}).Count(&totalCount)

	// Handle pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	if err := query.Find(&foods).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch foods"})
		return
	}

	response := serializers.SerializeFoods(foods, int(totalCount))
	c.JSON(http.StatusOK, response)
}

// GetFood returns a specific food
func GetFood(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var food models.Food
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&food).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
		return
	}

	response := serializers.SerializeFood(&food)
	c.JSON(http.StatusOK, response)
}

// CreateFoodRequest represents food creation request
type CreateFoodRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

// CreateFood creates a new food
func CreateFood(c *gin.Context) {
	var req CreateFoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	food, err := models.CreateFood(space, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create food"})
		return
	}

	response := serializers.SerializeFood(food)
	c.JSON(http.StatusCreated, response)
}

// UpdateFoodRequest represents food update request
type UpdateFoodRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdateFood updates a food
func UpdateFood(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID"})
		return
	}

	var req UpdateFoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var food models.Food
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&food).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
		return
	}

	// Update fields
	if req.Name != nil {
		food.Name = *req.Name
	}
	if req.Description != nil {
		food.Description = *req.Description
	}

	if err := models.DB.Save(&food).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update food"})
		return
	}

	response := serializers.SerializeFood(&food)
	c.JSON(http.StatusOK, response)
}

// DeleteFood deletes a food
func DeleteFood(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var food models.Food
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&food).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
		return
	}

	if err := models.DB.Delete(&food).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete food"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
