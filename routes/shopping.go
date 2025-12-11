package routes

import (
	"godoor/models"
	"godoor/serializers"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetShoppingListEntries returns list of shopping list entries with filtering and pagination
func GetShoppingListEntries(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)
	user := c.MustGet("user").(*models.User)

	// Parse query parameters
	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")
	mealplan := c.Query("mealplan")
	updatedAfter := c.Query("updated_after")
	// autosync := c.Query("autosync") // TODO: implement autosync

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
	query := models.DB.Where("space_id = ?", space.ID)

	// Filter by user permissions (user or shared users)
	query = query.Where("created_by_id = ? OR created_by_id IN (SELECT user_id FROM user_spaces WHERE space_id = ?)",
		user.ID, space.ID)

	// Preload related data
	query = query.Preload("CreatedBy").
		Preload("Food").
		Preload("Unit").
		Preload("Ingredient")

	// Apply filters
	if mealplan != "" {
		query = query.Where("list_recipe_id IN (SELECT id FROM shopping_list_recipes WHERE meal_plan_id = ?)", mealplan)
	}

	// Filter by recent entries (only unchecked or recently completed)
	todayStart := time.Now().Truncate(24 * time.Hour)
	weekAgo := todayStart.AddDate(0, 0, -7)
	query = query.Where("checked = ? OR (checked = ? AND completed_at >= ?)", false, true, weekAgo)

	if updatedAfter != "" {
		if updatedTime, err := time.Parse(time.RFC3339, updatedAfter); err == nil {
			query = query.Where("updated_at >= ?", updatedTime)
		}
	}

	// Execute query
	var entries []models.ShoppingListEntry
	if err := query.Find(&entries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shopping list entries"})
		return
	}

	// Apply autosync filter if requested
	// TODO: Implement autosync filtering if needed

	// Limit results to prevent overload (like in Tandoor)
	if len(entries) > 1000 {
		entries = entries[:1000]
	}

	response := serializers.SerializeShoppingListEntries(entries, len(entries), page, pageSize)
	c.JSON(http.StatusOK, response)
}

// GetShoppingListEntry returns a specific shopping list entry
func GetShoppingListEntry(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shopping list entry ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	user := c.MustGet("user").(*models.User)

	var entry models.ShoppingListEntry
	query := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id))
	query = query.Where("created_by_id = ? OR created_by_id IN (SELECT user_id FROM user_spaces WHERE space_id = ?)",
		user.ID, space.ID)

	if err := query.Preload("CreatedBy").
		Preload("Food").
		Preload("Food.ShoppingLists").
		Preload("ShoppingLists").
		Preload("Unit").
		Preload("ListRecipe").
		Preload("ListRecipe.Recipe").
		Preload("Ingredient").
		First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shopping list entry not found"})
		return
	}

	response := serializers.SerializeShoppingListEntry(&entry)
	c.JSON(http.StatusOK, response)
}

// CreateShoppingListEntry creates a new shopping list entry
func CreateShoppingListEntry(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	user := c.MustGet("user").(*models.User)

	// Validate required fields
	if foodID, exists := req["food_id"]; !exists || foodID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "food_id is required"})
		return
	}
	if foodID, ok := req["food_id"].(float64); !ok || foodID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "food_id must be a positive integer"})
		return
	}

	entry := models.ShoppingListEntry{
		SpaceID:     space.ID,
		Space:       *space,
		CreatedByID: user.ID,
		CreatedBy:   *user,
	}

	// Parse basic fields
	if amount, ok := req["amount"].(float64); ok {
		entry.Amount = amount
	}
	if order, ok := req["order"].(float64); ok {
		entry.Order = int(order)
	}
	if checked, ok := req["checked"].(bool); ok {
		entry.Checked = checked
		if checked {
			now := time.Now()
			entry.CompletedAt = &now
		}
	}

	// Parse food
	if foodID, ok := req["food_id"].(float64); ok {
		entry.FoodID = uint(foodID)
		// Load and validate food exists
		var food models.Food
		if err := models.DB.First(&food, uint(foodID)).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Food with the specified food_id does not exist"})
			return
		}
		entry.Food = food
	}

	// Parse unit
	if unitID, ok := req["unit_id"].(float64); ok && unitID > 0 {
		unitIDUint := uint(unitID)
		entry.UnitID = &unitIDUint
		// Load unit
		var unit models.Unit
		if err := models.DB.First(&unit, unitIDUint).Error; err == nil {
			entry.Unit = &unit
		}
	}

	// Parse ingredient
	if ingredientID, ok := req["ingredient_id"].(float64); ok && ingredientID > 0 {
		ingredientIDUint := uint(ingredientID)
		entry.IngredientID = &ingredientIDUint
		// Load ingredient
		var ingredient models.Ingredient
		if err := models.DB.First(&ingredient, ingredientIDUint).Error; err == nil {
			entry.Ingredient = &ingredient
		}
	}

	// Parse mealplan_id for creating ShoppingListRecipe
	if mealplanID, ok := req["mealplan_id"].(float64); ok && mealplanID > 0 {
		mealplanIDUint := uint(mealplanID)
		// Check if ShoppingListRecipe already exists for this mealplan
		var existingRecipe models.ShoppingListRecipe
		if err := models.DB.Where("meal_plan_id = ? AND space_id = ?", mealplanIDUint, space.ID).First(&existingRecipe).Error; err != nil {
			// Create new ShoppingListRecipe
			newRecipe := models.ShoppingListRecipe{
				MealPlanID:  &mealplanIDUint,
				SpaceID:     space.ID,
				Space:       *space,
				CreatedByID: user.ID,
				CreatedBy:   *user,
			}
			if err := models.DB.Create(&newRecipe).Error; err == nil {
				entry.ListRecipeID = &newRecipe.ID
				entry.ListRecipe = &newRecipe
			}
		} else {
			entry.ListRecipeID = &existingRecipe.ID
			entry.ListRecipe = &existingRecipe
		}
	}

	if err := models.DB.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shopping list entry"})
		return
	}

	// Load shopping lists for food and assign to entry
	if entry.Food.ID != 0 {
		var food models.Food
		if err := models.DB.Preload("ShoppingLists").First(&food, entry.FoodID).Error; err == nil {
			// Add entry to all shopping lists of the food
			for _, list := range food.ShoppingLists {
				models.DB.Exec("INSERT OR IGNORE INTO shopping_list_entry_lists (shopping_list_entry_id, shopping_list_id) VALUES (?, ?)",
					entry.ID, list.ID)
			}
		}
	}

	response := serializers.SerializeShoppingListEntry(&entry)
	c.JSON(http.StatusCreated, response)
}

// UpdateShoppingListEntry updates a shopping list entry
func UpdateShoppingListEntry(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shopping list entry ID"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	user := c.MustGet("user").(*models.User)

	var entry models.ShoppingListEntry
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shopping list entry not found"})
		return
	}

	// Check permissions
	if entry.CreatedByID != user.ID {
		// TODO: Check if user has access through sharing
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update fields
	if amount, ok := req["amount"].(float64); ok {
		entry.Amount = amount
	}
	if order, ok := req["order"].(float64); ok {
		entry.Order = int(order)
	}
	if checked, ok := req["checked"].(bool); ok {
		entry.Checked = checked
		if checked {
			now := time.Now()
			entry.CompletedAt = &now
		} else {
			entry.CompletedAt = nil
		}
		entry.UpdatedAt = time.Now()
	}

	// Handle onhand updates based on user preferences
	var userPref models.UserPreference
	if err := models.DB.Where("user_id = ?", user.ID).First(&userPref).Error; err == nil {
		if userPref.ShoppingAddOnhand {
			var food models.Food
			if err := models.DB.First(&food, entry.FoodID).Error; err == nil {
				if entry.Checked {
					// Add to onhand for user and shared users
					// TODO: Implement onhand_users relationship
				} else {
					// Remove from onhand
					// TODO: Implement onhand_users relationship
				}
			}
		}
	}

	if err := models.DB.Save(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shopping list entry"})
		return
	}

	response := serializers.SerializeShoppingListEntry(&entry)
	c.JSON(http.StatusOK, response)
}

// DeleteShoppingListEntry deletes a shopping list entry
func DeleteShoppingListEntry(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shopping list entry ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	user := c.MustGet("user").(*models.User)

	var entry models.ShoppingListEntry
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shopping list entry not found"})
		return
	}

	// Check permissions
	if entry.CreatedByID != user.ID {
		// TODO: Check if user has access through sharing
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := models.DB.Delete(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete shopping list entry"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// BulkUpdateShoppingListEntries updates multiple shopping list entries
func BulkUpdateShoppingListEntries(c *gin.Context) {
	var req struct {
		IDs     []uint `json:"ids" binding:"required"`
		Checked bool   `json:"checked" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	user := c.MustGet("user").(*models.User)

	// Update entries
	now := time.Now()
	updateData := map[string]interface{}{
		"checked":    req.Checked,
		"updated_at": now,
	}

	if req.Checked {
		updateData["completed_at"] = now
	} else {
		updateData["completed_at"] = nil
	}

	result := models.DB.Model(&models.ShoppingListEntry{}).
		Where("space_id = ? AND id IN (?) AND (created_by_id = ? OR created_by_id IN (SELECT user_id FROM user_spaces WHERE space_id = ?))",
			space.ID, req.IDs, user.ID, space.ID).
		Updates(updateData)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shopping list entries"})
		return
	}

	// Handle onhand updates if needed
	var userPref models.UserPreference
	if err := models.DB.Where("user_id = ?", user.ID).First(&userPref).Error; err == nil {
		if userPref.ShoppingAddOnhand {
			// TODO: Implement onhand updates for bulk operations
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"updated":   result.RowsAffected,
		"timestamp": now.Format(time.RFC3339),
	})
}
