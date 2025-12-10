package routes

import (
	"fmt"
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
	query := models.DB.Where("space_id = ?", space.ID).Preload("CreatedBy").Preload("Steps")

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
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).Preload("CreatedBy").Preload("Steps").First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	response := serializers.SerializeRecipe(&recipe)
	c.JSON(http.StatusOK, response)
}

// StepCreate represents a step in create request
type StepCreate struct {
	Name        string `json:"name"`
	Instruction string `json:"instruction"`
	Time        int    `json:"time"`
	Order       int    `json:"order"`
}

// CreateRecipeRequest represents recipe creation request
type CreateRecipeRequest struct {
	Name        string       `json:"name" binding:"required"`
	Description *string      `json:"description,omitempty"`
	Servings    int          `json:"servings,omitempty"`
	Steps       []StepCreate `json:"steps"`
}

// CreateRecipe creates a new recipe
func CreateRecipe(c *gin.Context) {
	var req CreateRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Debug logging
	fmt.Printf("DEBUG: Recipe name: %s, Steps count: %d\n", req.Name, len(req.Steps))

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

	// Create steps if provided
	if len(req.Steps) > 0 {
		for _, stepData := range req.Steps {
			step := models.Step{
				Name:        stepData.Name,
				Instruction: stepData.Instruction,
				Time:        stepData.Time,
				Order:       stepData.Order,
				RecipeID:    recipe.ID,
				Recipe:      *recipe,
				SpaceID:     space.ID,
				Space:       *space,
			}

			if err := models.DB.Create(&step).Error; err != nil {
				continue // Skip failed steps
			}
		}
	}

	if err := models.DB.Preload("CreatedBy").Preload("Steps").First(recipe, recipe.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload recipe"})
		return
	}

	response := serializers.SerializeRecipe(recipe)
	c.JSON(http.StatusCreated, response)
}

// UpdateRecipeRequest represents recipe update request
type UpdateRecipeRequest struct {
	Name        *string                   `json:"name,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Servings    *int                      `json:"servings,omitempty"`
	Steps       []map[string]interface{}  `json:"steps,omitempty"`       // TODO: implement proper step handling
	Properties  []map[string]interface{}  `json:"properties,omitempty"`  // TODO: implement proper property handling
}

// StepUpdate represents a step in update request
type StepUpdate struct {
	Name        string                   `json:"name"`
	Instruction string                   `json:"instruction"`
	Time        int                      `json:"time"`
	Ingredients []map[string]interface{} `json:"ingredients"`
	StepRecipe  *int                     `json:"step_recipe"`
	Order       int                      `json:"order"`
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

	// Handle steps - delete existing and create new ones
	var responseSteps []map[string]interface{}

	if req.Steps != nil && len(req.Steps) > 0 {
		// Delete existing steps for this recipe
		models.DB.Where("recipe_id = ?", recipe.ID).Delete(&models.Step{})

		// Create new steps
		responseSteps = make([]map[string]interface{}, len(req.Steps))
		for i, stepData := range req.Steps {
			step := models.Step{
				Name:        getStringFromMap(stepData, "name"),
				Instruction: getStringFromMap(stepData, "instruction"),
				Time:        getIntFromMap(stepData, "time"),
				Order:       getIntFromMap(stepData, "order"),
				SpaceID:     space.ID,
				Space:       *space,
				RecipeID:    recipe.ID,
				Recipe:      recipe,
			}

			if err := models.DB.Create(&step).Error; err != nil {
				continue
			}

			// Create response step data with real ID
			responseSteps[i] = map[string]interface{}{
				"id":                   step.ID,
				"name":                 step.Name,
				"instruction":          step.Instruction,
				"time":                 step.Time,
				"order":                step.Order,
				"show_as_header":       step.ShowAsHeader,
				"show_ingredients_table": step.ShowIngredientsTable,
				"ingredients":          []interface{}{}, // TODO: implement ingredients
				"instructions_markdown": step.Instruction, // TODO: implement markdown
				"file":                 nil,
				"step_recipe":          nil,
				"step_recipe_data":     nil,
				"numrecipe":           0,
			}
		}
	} else {
		// If no steps provided in request, get existing steps for the recipe
		var existingSteps []models.Step
		if err := models.DB.Where("recipe_id = ?", recipe.ID).Find(&existingSteps).Error; err == nil {
			responseSteps = make([]map[string]interface{}, len(existingSteps))
			for i, step := range existingSteps {
				responseSteps[i] = map[string]interface{}{
					"id":                   step.ID,
					"name":                 step.Name,
					"instruction":          step.Instruction,
					"time":                 step.Time,
					"order":                step.Order,
					"show_as_header":       step.ShowAsHeader,
					"show_ingredients_table": step.ShowIngredientsTable,
					"ingredients":          []interface{}{}, // TODO: implement ingredients
					"instructions_markdown": step.Instruction, // TODO: implement markdown
					"file":                 nil,
					"step_recipe":          nil,
					"step_recipe_data":     nil,
					"numrecipe":           0,
				}
			}
		} else {
			// If no existing steps, return empty array
			responseSteps = []map[string]interface{}{}
		}
	}

	if err := models.DB.Preload("CreatedBy").First(&recipe, recipe.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload recipe"})
		return
	}

	// Create response manually with proper steps
	response := map[string]interface{}{
		"id":                    recipe.ID,
		"name":                  recipe.Name,
		"description":           recipe.Description,
		"image":                 nil,
		"keywords":              []interface{}{},
		"steps":                 responseSteps,
		"working_time":          recipe.WorkingTime,
		"waiting_time":          recipe.WaitingTime,
		"created_by":            serializers.SerializeUser(&recipe.CreatedBy),
		"created_at":            recipe.CreatedAt,
		"updated_at":            recipe.UpdatedAt,
		"source_url":            "",
		"internal":              recipe.Internal,
		"show_ingredient_overview": true,
		"nutrition":             nil,
		"properties":            []interface{}{},
		"food_properties":       map[string]interface{}{},
		"servings":              recipe.Servings,
		"file_path":             "",
		"servings_text":         recipe.ServingsText,
		"rating":                recipe.Rating,
		"last_cooked":           nil,
		"private":               recipe.Private,
		"shared":                []interface{}{},
	}

	c.JSON(http.StatusOK, response)
}

// Helper functions to extract values from map
func getStringFromMap(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getIntFromMap(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
	}
	return 0
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
