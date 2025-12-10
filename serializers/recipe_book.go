package serializers

import (
	"godoor/models"
	"time"
)

// RecipeBookSerializer matching Tandoor RecipeBookSerializer
type RecipeBookSerializer struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Shared      []UserIDSerializer `json:"shared"`
	CreatedBy   UserIDSerializer `json:"created_by"`
	Filter      interface{}      `json:"filter"` // TODO: implement CustomFilter
	Order       int              `json:"order"`
	Count       int              `json:"count"` // TODO: implement recipe count
}

// SerializeRecipeBook converts RecipeBook model to serializer
func SerializeRecipeBook(book *models.RecipeBook) RecipeBookSerializer {
	return RecipeBookSerializer{
		ID:          book.ID,
		Name:        book.Name,
		Description: book.Description,
		Shared:      []UserIDSerializer{}, // TODO: implement shared users
		CreatedBy:   SerializeUserID(&book.CreatedBy),
		Filter:      nil, // TODO: implement filter
		Order:       book.Order,
		Count:       0,   // TODO: implement recipe count
	}
}

// RecipeBookListResponse represents paginated recipe book response matching Django REST Framework
type RecipeBookListResponse struct {
	Count     int                     `json:"count"`
	Next      *string                 `json:"next"`
	Previous  *string                 `json:"previous"`
	Results   []RecipeBookSerializer   `json:"results"`
	Timestamp string                  `json:"timestamp"`
}

// SerializeRecipeBooks converts slice of RecipeBook models to paginated response
func SerializeRecipeBooks(books []models.RecipeBook) RecipeBookListResponse {
	result := make([]RecipeBookSerializer, len(books))
	for i, book := range books {
		result[i] = SerializeRecipeBook(&book)
	}
	return RecipeBookListResponse{
		Count:     len(books),
		Next:      nil, // TODO: implement pagination
		Previous:  nil, // TODO: implement pagination
		Results:   result,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
