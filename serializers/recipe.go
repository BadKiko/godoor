package serializers

import (
	"fmt"
	"godoor/config"
	"godoor/models"
	"time"
)

// getImageURL returns full image URL or nil if no image
func getImageURL(imagePath *string) interface{} {
	if imagePath != nil && *imagePath != "" {
		// If it's already a full URL (starts with http), return as is
		if len(*imagePath) > 4 && (*imagePath)[:4] == "http" {
			return *imagePath
		}
		// Otherwise, prepend MEDIA_URL
		fullURL := config.MediaURL + *imagePath
		return fullURL
	}
	return nil
}

// RecipeOverviewSerializer matching Tandoor RecipeOverviewSerializer
type RecipeOverviewSerializer struct {
	ID           uint                     `json:"id"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description,omitempty"`
	Image        interface{}              `json:"image"`
	Keywords     []KeywordLabelSerializer `json:"keywords"`
	New          bool                     `json:"new"`
	Recent       string                   `json:"recent"`
	Rating       *float64                 `json:"rating,omitempty"`
	LastCooked   *time.Time               `json:"last_cooked,omitempty"`
	CreatedBy    UserSerializer           `json:"created_by"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	Internal     bool                     `json:"internal"`
	Private      bool                     `json:"private"`
	Servings     int                      `json:"servings"`
	ServingsText string                   `json:"servings_text"`
	WorkingTime  int                      `json:"working_time"`
	WaitingTime  int                      `json:"waiting_time"`
}

// SerializeRecipeOverview converts Recipe model to overview serializer
func SerializeRecipeOverview(recipe *models.Recipe) RecipeOverviewSerializer {
	description := ""
	if recipe.Description != nil {
		description = *recipe.Description
	}

	// Load keywords for recipe - they should be preloaded by the query
	var keywords []KeywordLabelSerializer
	for _, keyword := range recipe.Keywords {
		keywords = append(keywords, SerializeKeywordLabel(&keyword))
	}

	// TODO: implement new recipe logic (recipe created within last 7 days?)
	isNew := time.Since(recipe.CreatedAt).Hours() < 24*7

	// TODO: implement recent logic
	recent := ""

	return RecipeOverviewSerializer{
		ID:           recipe.ID,
		Name:         recipe.Name,
		Description:  description,
		Image:        getImageURL(recipe.Image),
		Keywords:     keywords,
		New:          isNew,
		Recent:       recent,
		Rating:       recipe.Rating,
		LastCooked:   nil, // TODO: implement last cooked from cook logs
		CreatedBy:    SerializeUser(&recipe.CreatedBy),
		CreatedAt:    recipe.CreatedAt,
		UpdatedAt:    recipe.UpdatedAt,
		Internal:     recipe.Internal,
		Private:      recipe.Private,
		Servings:     recipe.Servings,
		ServingsText: recipe.ServingsText,
		WorkingTime:  recipe.WorkingTime,
		WaitingTime:  recipe.WaitingTime,
	}
}

// RecipeListResponse represents paginated recipe response matching Django REST Framework
type RecipeListResponse struct {
	Count     int                        `json:"count"`
	Next      *string                    `json:"next"`
	Previous  *string                    `json:"previous"`
	Results   []RecipeOverviewSerializer `json:"results"`
	Timestamp string                     `json:"timestamp"`
}

// SerializeRecipes converts slice of Recipe models to paginated response
func SerializeRecipes(recipes []models.Recipe, totalCount int, page int, pageSize int) RecipeListResponse {
	result := make([]RecipeOverviewSerializer, len(recipes))
	for i, recipe := range recipes {
		result[i] = SerializeRecipeOverview(&recipe)
	}

	// Calculate pagination URLs
	var next, previous *string

	// Calculate if there are more pages
	totalPages := (totalCount + pageSize - 1) / pageSize // Ceiling division

	if page < totalPages {
		nextURL := fmt.Sprintf("?page=%d&page_size=%d", page+1, pageSize)
		next = &nextURL
	}

	if page > 1 {
		prevURL := fmt.Sprintf("?page=%d&page_size=%d", page-1, pageSize)
		previous = &prevURL
	}

	return RecipeListResponse{
		Count:     totalCount,
		Next:      next,
		Previous:  previous,
		Results:   result,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
