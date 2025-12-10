package routes

import (
	"net/http"
	"strconv"
	"godoor/models"
	"godoor/serializers"
	"github.com/gin-gonic/gin"
)

// RecipeBookEntrySerializer represents the relationship between recipes and books
type RecipeBookEntrySerializer struct {
	ID       uint                           `json:"id"`
	Recipe   serializers.RecipeOverviewSerializer `json:"recipe"`
	Book     serializers.RecipeBookSerializer    `json:"book"`
	BookID   uint                           `json:"-"` // For creation
	RecipeID uint                           `json:"-"` // For creation
}

// RecipeBookEntryListResponse represents paginated recipe book entry response
type RecipeBookEntryListResponse struct {
	Count     int                          `json:"count"`
	Next      *string                      `json:"next"`
	Previous  *string                      `json:"previous"`
	Results   []RecipeBookEntrySerializer  `json:"results"`
	Timestamp string                       `json:"timestamp"`
}

// GetRecipeBookEntries returns list of recipe book entries for current user
func GetRecipeBookEntries(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)

	// Handle book filter
	bookIDStr := c.Query("book")
	var bookEntries []models.RecipeBookEntry

	// First get book IDs for this space
	var bookIDs []uint
	if err := models.DB.Model(&models.RecipeBook{}).Where("space_id = ?", space.ID).Pluck("id", &bookIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	if len(bookIDs) == 0 {
		// No books in this space, return empty result
		response := RecipeBookEntryListResponse{
			Count:     0,
			Next:      nil,
			Previous:  nil,
			Results:   []RecipeBookEntrySerializer{},
			Timestamp: "2025-12-10T14:15:00+03:00",
		}
		c.JSON(http.StatusOK, response)
		return
	}

	query := models.DB.Where("book_id IN ?", bookIDs)

	if bookIDStr != "" {
		bookID, err := strconv.Atoi(bookIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
			return
		}
		// Verify the book belongs to this space
		found := false
		for _, id := range bookIDs {
			if id == uint(bookID) {
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		query = query.Where("book_id = ?", uint(bookID))
	}

	// Handle pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	var totalCount int64
	query.Model(&models.RecipeBookEntry{}).Count(&totalCount)

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	if err := query.Preload("Recipe").Preload("Book").Find(&bookEntries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipe book entries"})
		return
	}

	// Convert to serializers
	results := make([]RecipeBookEntrySerializer, len(bookEntries))
	for i, entry := range bookEntries {
		results[i] = RecipeBookEntrySerializer{
			ID:       entry.ID,
			Recipe:   serializers.SerializeRecipeOverview(&entry.Recipe),
			Book:     serializers.SerializeRecipeBook(&entry.Book),
			BookID:   entry.BookID,
			RecipeID: entry.RecipeID,
		}
	}

	response := RecipeBookEntryListResponse{
		Count:     int(totalCount),
		Next:      nil, // TODO: implement pagination
		Previous:  nil, // TODO: implement pagination
		Results:   results,
		Timestamp: "2025-12-10T14:15:00+03:00", // TODO: use actual timestamp
	}

	c.JSON(http.StatusOK, response)
}

// CreateRecipeBookEntryRequest represents recipe book entry creation request
type CreateRecipeBookEntryRequest struct {
	BookID   uint `json:"book" binding:"required"`
	RecipeID uint `json:"recipe" binding:"required"`
}

// CreateRecipeBookEntry creates a new recipe book entry
func CreateRecipeBookEntry(c *gin.Context) {
	var req CreateRecipeBookEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	// Verify book belongs to space
	var book models.RecipeBook
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, req.BookID).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Verify recipe belongs to space
	var recipe models.Recipe
	if err := models.DB.Where("space_id = ? AND id = ?", space.ID, req.RecipeID).First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	// Check if entry already exists
	var existingEntry models.RecipeBookEntry
	if err := models.DB.Where("book_id = ? AND recipe_id = ?", req.BookID, req.RecipeID).First(&existingEntry).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Recipe already in book"})
		return
	}

	entry := models.RecipeBookEntry{
		BookID:   req.BookID,
		Book:     book,
		RecipeID: req.RecipeID,
		Recipe:   recipe,
	}

	if err := models.DB.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe book entry"})
		return
	}

	response := RecipeBookEntrySerializer{
		ID:       entry.ID,
		Recipe:   serializers.SerializeRecipeOverview(&recipe),
		Book:     serializers.SerializeRecipeBook(&book),
		BookID:   entry.BookID,
		RecipeID: entry.RecipeID,
	}

	c.JSON(http.StatusCreated, response)
}

// DeleteRecipeBookEntry deletes a recipe book entry
func DeleteRecipeBookEntry(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe book entry ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	// Delete only if the book belongs to the space
	result := models.DB.Joins("Book").Where("recipe_book_entries.id = ? AND recipe_books.space_id = ?", uint(id), space.ID).Delete(&models.RecipeBookEntry{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe book entry"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe book entry not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
