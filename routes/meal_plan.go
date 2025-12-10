package routes

import (
	"net/http"
	"strconv"
	"time"
	"godoor/models"
	"godoor/serializers"
	"github.com/gin-gonic/gin"
)

// GetMealTypes returns list of meal types for current user
func GetMealTypes(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)

	var mealTypes []models.MealType
	if err := models.DB.Where("space_id = ?", space.ID).Find(&mealTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch meal types"})
		return
	}

	response := serializers.SerializeMealTypes(mealTypes)
	c.JSON(http.StatusOK, response)
}

// CreateMealTypeRequest represents meal type creation request
type CreateMealTypeRequest struct {
	Name   string  `json:"name" binding:"required"`
	Order  int     `json:"order,omitempty"`
	Color  *string `json:"color,omitempty"`
	Time   *string `json:"time,omitempty"`
	Default bool   `json:"default,omitempty"`
}

// CreateMealType creates a new meal type
func CreateMealType(c *gin.Context) {
	var req CreateMealTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user := c.MustGet("user").(*models.User)
	space := c.MustGet("space").(*models.Space)

	mealType, err := models.CreateMealType(user, space, req.Name, req.Order, req.Color, req.Time, req.Default)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal type"})
		return
	}

	response := serializers.SerializeMealType(mealType)
	c.JSON(http.StatusCreated, response)
}

// GetMealPlans returns list of meal plans with date filtering
func GetMealPlans(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)

	// Parse date filters
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	var mealPlans []models.MealPlan
	query := models.DB.Where("space_id = ?", space.ID)

	if fromDateStr != "" {
		fromDate, err := time.Parse("2006-01-02", fromDateStr)
		if err == nil {
			query = query.Where("from_date >= ?", fromDate)
		}
	}

	if toDateStr != "" {
		toDate, err := time.Parse("2006-01-02", toDateStr)
		if err == nil {
			query = query.Where("to_date <= ?", toDate)
		}
	}

	if err := query.Preload("MealType").Preload("Recipe").Preload("CreatedBy").Find(&mealPlans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch meal plans"})
		return
	}

	response := serializers.SerializeMealPlans(mealPlans)
	c.JSON(http.StatusOK, response)
}

// GetMealPlan returns a specific meal plan
func GetMealPlan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal plan ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var mealPlan models.MealPlan
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).
		Preload("MealType").Preload("Recipe").Preload("CreatedBy").First(&mealPlan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meal plan not found"})
		return
	}

	response := serializers.SerializeMealPlan(&mealPlan)
	c.JSON(http.StatusOK, response)
}

// CreateMealPlanRequest represents meal plan creation request
type CreateMealPlanRequest struct {
	Title      string  `json:"title,omitempty"`
	RecipeID   *uint   `json:"recipe,omitempty"`
	Servings   float64 `json:"servings,omitempty"`
	MealTypeID uint    `json:"meal_type" binding:"required"`
	Note       string  `json:"note,omitempty"`
	FromDate   string  `json:"from_date" binding:"required"`
	ToDate     string  `json:"to_date,omitempty"`
}

// CreateMealPlan creates a new meal plan
func CreateMealPlan(c *gin.Context) {
	var req CreateMealPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user := c.MustGet("user").(*models.User)
	space := c.MustGet("space").(*models.Space)

	// Find meal type
	var mealType models.MealType
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, req.MealTypeID).First(&mealType).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal type"})
		return
	}

	// Parse dates
	fromDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.FromDate)
	if err != nil {
		fromDate, err = time.Parse("2006-01-02", req.FromDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from_date format"})
			return
		}
	}

	toDate := fromDate
	if req.ToDate != "" {
		toDate, err = time.Parse("2006-01-02T15:04:05Z07:00", req.ToDate)
		if err != nil {
			toDate, err = time.Parse("2006-01-02", req.ToDate)
			if err != nil {
				toDate = fromDate
			}
		}
	}

	servings := req.Servings
	if servings <= 0 {
		servings = 1
	}

	mealPlan, err := models.CreateMealPlan(user, space, &mealType, fromDate, toDate, req.Title, req.Note, req.RecipeID, servings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal plan"})
		return
	}

	if err := models.DB.Preload("MealType").Preload("Recipe").Preload("CreatedBy").First(mealPlan, mealPlan.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload meal plan"})
		return
	}

	response := serializers.SerializeMealPlan(mealPlan)
	c.JSON(http.StatusCreated, response)
}

// UpdateMealPlanRequest represents meal plan update request
type UpdateMealPlanRequest struct {
	Title      *string  `json:"title,omitempty"`
	RecipeID   *uint    `json:"recipe,omitempty"`
	Servings   *float64 `json:"servings,omitempty"`
	MealTypeID *uint    `json:"meal_type,omitempty"`
	Note       *string  `json:"note,omitempty"`
	FromDate   *string  `json:"from_date,omitempty"`
	ToDate     *string  `json:"to_date,omitempty"`
}

// UpdateMealPlan updates a meal plan
func UpdateMealPlan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal plan ID"})
		return
	}

	var req UpdateMealPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	currentUser := c.MustGet("user").(*models.User)

	var mealPlan models.MealPlan
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&mealPlan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meal plan not found"})
		return
	}

	// Check permissions (only creator can update)
	if mealPlan.CreatedByID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Update fields
	if req.Title != nil {
		mealPlan.Title = *req.Title
	}
	if req.RecipeID != nil {
		mealPlan.RecipeID = req.RecipeID
	}
	if req.Servings != nil && *req.Servings > 0 {
		mealPlan.Servings = *req.Servings
	}
	if req.MealTypeID != nil {
		var mealType models.MealType
		if err := models.DB.Where("space_id = ? AND id = ?", space.ID, *req.MealTypeID).First(&mealType).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal type"})
			return
		}
		mealPlan.MealTypeID = *req.MealTypeID
	}
	if req.Note != nil {
		mealPlan.Note = *req.Note
	}
	if req.FromDate != nil {
		if fromDate, err := time.Parse("2006-01-02T15:04:05Z07:00", *req.FromDate); err == nil {
			mealPlan.FromDate = fromDate
		} else if fromDate, err := time.Parse("2006-01-02", *req.FromDate); err == nil {
			mealPlan.FromDate = fromDate
		}
	}
	if req.ToDate != nil {
		if toDate, err := time.Parse("2006-01-02T15:04:05Z07:00", *req.ToDate); err == nil {
			mealPlan.ToDate = toDate
		} else if toDate, err := time.Parse("2006-01-02", *req.ToDate); err == nil {
			mealPlan.ToDate = toDate
		}
	}

	if err := models.DB.Save(&mealPlan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meal plan"})
		return
	}

	if err := models.DB.Preload("MealType").Preload("Recipe").Preload("CreatedBy").First(&mealPlan, mealPlan.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload meal plan"})
		return
	}

	response := serializers.SerializeMealPlan(&mealPlan)
	c.JSON(http.StatusOK, response)
}

// DeleteMealPlan deletes a meal plan
func DeleteMealPlan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal plan ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	currentUser := c.MustGet("user").(*models.User)

	var mealPlan models.MealPlan
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, uint(id)).First(&mealPlan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meal plan not found"})
		return
	}

	// Check permissions (only creator can delete)
	if mealPlan.CreatedByID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := models.DB.Delete(&mealPlan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meal plan"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
