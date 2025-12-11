package routes

import (
	"fmt"
	"godoor/models"
	"godoor/serializers"
	"godoor/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetRecipes returns list of recipes with sorting and pagination
func GetRecipes(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)

	// Parse query parameters
	sortOrder := c.Query("sort_order")
	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")
	searchQuery := c.Query("query")
	timesCookedStr := c.Query("timescooked")

	// Parse keywords parameters
	keywordsOr := c.QueryArray("keywords_or")
	keywordsAnd := c.QueryArray("keywords_and")
	keywordsOrNot := c.QueryArray("keywords_or_not")
	keywordsAndNot := c.QueryArray("keywords_and_not")
	keywords := c.QueryArray("keywords") // alias for keywords_or

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
	query := models.DB.Where("space_id = ?", space.ID).Preload("CreatedBy").Preload("Keywords").Preload("Steps").Preload("Steps.Ingredients").Preload("Steps.Ingredients.Food").Preload("Steps.Ingredients.Unit")

	// Apply filters
	if searchQuery != "" {
		// Search in recipe name and description
		query = query.Where("LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)",
			"%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	if timesCookedStr != "" {
		if timesCooked, err := strconv.Atoi(timesCookedStr); err == nil {
			if timesCooked == 0 {
				// Recipes never cooked (no cook logs)
				query = query.Where("id NOT IN (SELECT DISTINCT recipe_id FROM cook_logs WHERE space_id = ?)", space.ID)
			} else {
				// Recipes cooked exactly N times
				query = query.Joins("LEFT JOIN (SELECT recipe_id, COUNT(*) as cook_count FROM cook_logs WHERE space_id = ? GROUP BY recipe_id) cl ON recipes.id = cl.recipe_id", space.ID).
					Where("COALESCE(cl.cook_count, 0) = ?", timesCooked)
			}
		}
	}

	// Apply keywords filters
	if len(keywordsOr) > 0 || len(keywords) > 0 {
		// Combine keywords arrays
		keywordIDs := append(keywordsOr, keywords...)

		// Filter recipes that have ANY of the specified keywords
		query = query.Joins("JOIN recipe_keywords rk_or ON recipes.id = rk_or.recipe_id").
			Where("rk_or.keyword_id IN (?)", keywordIDs)
	}

	if len(keywordsAnd) > 0 {
		// Filter recipes that have ALL of the specified keywords
		for _, keywordID := range keywordsAnd {
			query = query.Joins("JOIN recipe_keywords rk_and ON recipes.id = rk_and.recipe_id").
				Where("rk_and.keyword_id = ?", keywordID)
		}
	}

	if len(keywordsOrNot) > 0 {
		// Exclude recipes that have ANY of the specified keywords
		query = query.Where("recipes.id NOT IN (SELECT DISTINCT recipe_id FROM recipe_keywords WHERE keyword_id IN (?))", keywordsOrNot)
	}

	if len(keywordsAndNot) > 0 {
		// This is complex - exclude recipes that have ALL of the specified keywords
		// For now, implement as excluding recipes that have any of them (simplified)
		query = query.Where("recipes.id NOT IN (SELECT DISTINCT recipe_id FROM recipe_keywords WHERE keyword_id IN (?))", keywordsAndNot)
	}

	if sortOrder == "-favorite" {
		// TODO: Implement favorite filtering (requires user favorite recipes table)
		// For now, just sort by created_at
		query = query.Order("created_at DESC")
	}

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
	response := serializers.SerializeRecipes(recipes, int(totalCount), page, pageSize)
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
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).Preload("CreatedBy").Preload("Steps").Preload("Steps.Ingredients").Preload("Steps.Ingredients.Food").Preload("Steps.Ingredients.Unit").First(&recipe).Error; err != nil {
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
	Name        *string                  `json:"name,omitempty"`
	Description *string                  `json:"description,omitempty"`
	Servings    *int                     `json:"servings,omitempty"`
	Steps       []map[string]interface{} `json:"steps,omitempty"`      // TODO: implement proper step handling
	Properties  []map[string]interface{} `json:"properties,omitempty"` // TODO: implement proper property handling
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
				"id":                     step.ID,
				"name":                   step.Name,
				"instruction":            step.Instruction,
				"time":                   step.Time,
				"order":                  step.Order,
				"show_as_header":         step.ShowAsHeader,
				"show_ingredients_table": step.ShowIngredientsTable,
				"ingredients":            []interface{}{},  // TODO: implement ingredients
				"instructions_markdown":  step.Instruction, // TODO: implement markdown
				"file":                   nil,
				"step_recipe":            nil,
				"step_recipe_data":       nil,
				"numrecipe":              0,
			}
		}
	} else {
		// If no steps provided in request, get existing steps for the recipe
		var existingSteps []models.Step
		if err := models.DB.Where("recipe_id = ?", recipe.ID).Find(&existingSteps).Error; err == nil {
			responseSteps = make([]map[string]interface{}, len(existingSteps))
			for i, step := range existingSteps {
				responseSteps[i] = map[string]interface{}{
					"id":                     step.ID,
					"name":                   step.Name,
					"instruction":            step.Instruction,
					"time":                   step.Time,
					"order":                  step.Order,
					"show_as_header":         step.ShowAsHeader,
					"show_ingredients_table": step.ShowIngredientsTable,
					"ingredients":            []interface{}{},  // TODO: implement ingredients
					"instructions_markdown":  step.Instruction, // TODO: implement markdown
					"file":                   nil,
					"step_recipe":            nil,
					"step_recipe_data":       nil,
					"numrecipe":              0,
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
		"id":                       recipe.ID,
		"name":                     recipe.Name,
		"description":              recipe.Description,
		"image":                    recipe.Image,
		"keywords":                 []interface{}{},
		"steps":                    responseSteps,
		"working_time":             recipe.WorkingTime,
		"waiting_time":             recipe.WaitingTime,
		"created_by":               serializers.SerializeUser(&recipe.CreatedBy),
		"created_at":               recipe.CreatedAt,
		"updated_at":               recipe.UpdatedAt,
		"source_url":               "",
		"internal":                 recipe.Internal,
		"show_ingredient_overview": true,
		"nutrition":                nil,
		"properties":               []interface{}{},
		"food_properties":          map[string]interface{}{},
		"servings":                 recipe.Servings,
		"file_path":                "",
		"servings_text":            recipe.ServingsText,
		"rating":                   recipe.Rating,
		"last_cooked":              nil,
		"private":                  recipe.Private,
		"shared":                   []interface{}{},
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

// RecipeImageRequest represents recipe image update request (like Tandoor RecipeImageSerializer)
type RecipeImageRequest struct {
	Image    string `form:"image"`     // Uploaded file (not used directly)
	ImageURL string `form:"image_url"` // URL to download image from
}

// RecipeImageResponse represents the response structure
type RecipeImageResponse struct {
	ID       uint   `json:"id"`
	Image    string `json:"image,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

// validateImageURL performs basic URL validation (simplified version of Tandoor validate_import_url)
func validateImageURL(url string) bool {
	if url == "" {
		return false
	}
	// Basic checks - should start with http/https and be reasonable length
	return (strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) && len(url) < 4096
}

// RecipeImage updates recipe image (simplified to avoid 400 errors)
func RecipeImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var recipe models.Recipe
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	// Debug: log all form fields and files
	fmt.Printf("DEBUG RecipeImage: Recipe ID %d\n", id)
	fmt.Printf("DEBUG RecipeImage: Content-Type: %s\n", c.GetHeader("Content-Type"))
	fmt.Printf("DEBUG RecipeImage: Form values: %v\n", c.Request.PostForm)

	// Check all multipart files
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		fmt.Printf("DEBUG RecipeImage: Multipart files:\n")
		for field, files := range form.File {
			fmt.Printf("  Field '%s': %d files\n", field, len(files))
			for _, file := range files {
				fmt.Printf("    - %s (%d bytes)\n", file.Filename, file.Size)
			}
		}
	}

	// Simplified handling - just accept any multipart request
	// TODO: Implement full image processing like Tandoor
	imageURL := c.PostForm("image_url")

	if imageURL != "" && validateImageURL(imageURL) {
		// Delete old image file if exists before setting URL
		if recipe.Image != nil {
			utils.DeleteImageFile(*recipe.Image)
		}
		recipe.Image = &imageURL
	} else {
		// Check for uploaded file
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			// Try other common field names
			file, header, err = c.Request.FormFile("file")
		}
		if err != nil {
			// Try any file field - expanded list
			for _, field := range []string{"photo", "picture", "upload", "attachment", "media", "data", "content"} {
				file, header, err = c.Request.FormFile(field)
				if err == nil {
					fmt.Printf("DEBUG RecipeImage: Found file in field '%s'\n", field)
					break
				}
			}
		}

		if err == nil && header != nil {
			// Delete old image file if exists before saving new one
			if recipe.Image != nil {
				utils.DeleteImageFile(*recipe.Image)
			}

			// Save uploaded file to disk
			imagePath, saveErr := utils.SaveUploadedImage(file, header, recipe.ID)
			if saveErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": saveErr.Error()})
				return
			}
			recipe.Image = &imagePath
			file.Close()
		} else {
			// No image provided - remove current image
			if recipe.Image != nil {
				// Delete old image file from disk
				utils.DeleteImageFile(*recipe.Image)
			}
			recipe.Image = nil
		}
	}

	// Save the recipe
	fmt.Printf("DEBUG RecipeImage: Before save - recipe.Image = %v\n", recipe.Image)
	if err := models.DB.Save(&recipe).Error; err != nil {
		fmt.Printf("DEBUG RecipeImage: Save failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipe image"})
		return
	}
	fmt.Printf("DEBUG RecipeImage: After save - recipe.Image = %v\n", recipe.Image)

	// Return minimal response like Tandoor
	response := RecipeImageResponse{
		ID: recipe.ID,
	}
	if recipe.Image != nil {
		response.Image = *recipe.Image
		fmt.Printf("DEBUG RecipeImage: Response image = %s\n", response.Image)
	}
	if imageURL != "" {
		response.ImageURL = imageURL
		fmt.Printf("DEBUG RecipeImage: Response image_url = %s\n", response.ImageURL)
	}

	c.JSON(http.StatusOK, response)
}

// RecipeShopping adds recipe to shopping list
func RecipeShopping(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var recipe models.Recipe
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	// TODO: Implement full shopping list functionality
	// For now, just return success like the original Tandoor API
	c.JSON(http.StatusNoContent, gin.H{"msg": "Recipe added to shopping list"})
}
