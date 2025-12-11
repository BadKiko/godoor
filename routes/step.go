package routes

import (
	"godoor/models"
	"godoor/serializers"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetSteps returns list of steps with pagination and filtering
func GetSteps(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)

	// Parse query parameters
	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")
	recipeIDsStr := c.QueryArray("recipe")
	query := c.Query("query")

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
	queryBuilder := models.DB.Where("space_id = ?", space.ID)

	// Filter by recipe IDs if provided
	if len(recipeIDsStr) > 0 {
		var recipeIDs []uint
		for _, idStr := range recipeIDsStr {
			if id, err := strconv.Atoi(idStr); err == nil {
				recipeIDs = append(recipeIDs, uint(id))
			}
		}
		if len(recipeIDs) > 0 {
			queryBuilder = queryBuilder.Where("recipe_id IN ?", recipeIDs)
		}
	}

	// Filter by query (name or recipe name)
	if query != "" {
		queryBuilder = queryBuilder.Joins("JOIN recipes ON steps.recipe_id = recipes.id").
			Where("LOWER(steps.name) LIKE LOWER(?) OR LOWER(recipes.name) LIKE LOWER(?)",
				"%"+query+"%", "%"+query+"%")
	}

	// Get total count
	var totalCount int64
	queryBuilder.Model(&models.Step{}).Count(&totalCount)

	// Apply pagination
	offset := (page - 1) * pageSize
	queryBuilder = queryBuilder.Offset(offset).Limit(pageSize)

	// Execute query
	var steps []models.Step
	if err := queryBuilder.Preload("Recipe").Preload("Ingredients").Preload("Ingredients.Food").Preload("Ingredients.Unit").Find(&steps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch steps"})
		return
	}

	// Serialize response
	response := serializers.SerializeSteps(steps, int(totalCount))
	c.JSON(http.StatusOK, response)
}

// GetStep returns a specific step
func GetStep(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid step ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var step models.Step
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).Preload("Recipe").Preload("Ingredients").Preload("Ingredients.Food").Preload("Ingredients.Unit").First(&step).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Step not found"})
		return
	}

	response := serializers.SerializeStep(&step)
	c.JSON(http.StatusOK, response)
}

// CreateStepRequest represents step creation request
type CreateStepRequest struct {
	Name        string `json:"name"`
	Instruction string `json:"instruction"`
	Time        int    `json:"time"`
	Order       int    `json:"order"`
	RecipeID    uint   `json:"recipe" binding:"required"`
}

// CreateStep creates a new step
func CreateStep(c *gin.Context) {
	var req CreateStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	// Verify recipe exists and belongs to space
	var recipe models.Recipe
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, req.RecipeID).First(&recipe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Recipe not found"})
		return
	}

	step := models.Step{
		Name:        req.Name,
		Instruction: req.Instruction,
		Time:        req.Time,
		Order:       req.Order,
		SpaceID:     space.ID,
		Space:       *space,
		RecipeID:    req.RecipeID,
		Recipe:      recipe,
	}

	// Create step
	if err := models.DB.Create(&step).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create step"})
		return
	}

	if err := models.DB.Preload("Recipe").First(&step, step.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload step"})
		return
	}

	response := serializers.SerializeStep(&step)
	c.JSON(http.StatusCreated, response)
}

// UpdateStepRequest represents step update request
type UpdateStepRequest struct {
	Name        *string `json:"name,omitempty"`
	Instruction *string `json:"instruction,omitempty"`
	Time        *int    `json:"time,omitempty"`
	Order       *int    `json:"order,omitempty"`
}

// UpdateStep updates a step
func UpdateStep(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid step ID"})
		return
	}

	var req UpdateStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	currentUser := c.MustGet("user").(*models.User)

	var step models.Step
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&step).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Step not found"})
		return
	}

	// Check permissions - user can update steps in their recipes
	var recipe models.Recipe
	if err := models.DB.Where("id = ?", step.RecipeID).First(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find associated recipe"})
		return
	}

	if recipe.CreatedByID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Update fields
	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		step.Name = *req.Name
	}
	if req.Instruction != nil {
		step.Instruction = *req.Instruction
	}
	if req.Time != nil && *req.Time >= 0 {
		step.Time = *req.Time
	}
	if req.Order != nil && *req.Order >= 0 {
		step.Order = *req.Order
	}

	if err := models.DB.Save(&step).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update step"})
		return
	}

	if err := models.DB.Preload("Recipe").First(&step, step.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload step"})
		return
	}

	response := serializers.SerializeStep(&step)
	c.JSON(http.StatusOK, response)
}

// DeleteStep deletes a step
func DeleteStep(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid step ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	currentUser := c.MustGet("user").(*models.User)

	var step models.Step
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&step).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Step not found"})
		return
	}

	// Check permissions - user can delete steps in their recipes
	var recipe models.Recipe
	if err := models.DB.Where("id = ?", step.RecipeID).First(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find associated recipe"})
		return
	}

	if recipe.CreatedByID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := models.DB.Delete(&step).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete step"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
