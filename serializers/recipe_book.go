package serializers

import (
	"godoor/models"
)

// RecipeBookSerializer matching Tandoor RecipeBookSerializer
type RecipeBookSerializer struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Shared      []UserSerializer `json:"shared"`
	CreatedBy   UserSerializer `json:"created_by"`
	Filter      interface{}  `json:"filter"` // TODO: implement CustomFilter
	Order       int          `json:"order"`
	Count       int          `json:"count"` // TODO: implement recipe count
}

// SerializeRecipeBook converts RecipeBook model to serializer
func SerializeRecipeBook(book *models.RecipeBook) RecipeBookSerializer {
	return RecipeBookSerializer{
		ID:          book.ID,
		Name:        book.Name,
		Description: book.Description,
		Shared:      []UserSerializer{}, // TODO: implement shared users
		CreatedBy:   SerializeUser(&book.CreatedBy),
		Filter:      nil, // TODO: implement filter
		Order:       book.Order,
		Count:       0,   // TODO: implement recipe count
	}
}

// RecipeBookListResponse represents paginated recipe book response matching Django REST Framework
type RecipeBookListResponse struct {
	Count    int                     `json:"count"`
	Next     *string                 `json:"next"`
	Previous *string                 `json:"previous"`
	Results  []RecipeBookSerializer   `json:"results"`
}

// SerializeRecipeBooks converts slice of RecipeBook models to paginated response
func SerializeRecipeBooks(books []models.RecipeBook) RecipeBookListResponse {
	result := make([]RecipeBookSerializer, len(books))
	for i, book := range books {
		result[i] = SerializeRecipeBook(&book)
	}
	return RecipeBookListResponse{
		Count:    len(books),
		Next:     nil, // TODO: implement pagination
		Previous: nil, // TODO: implement pagination
		Results:  result,
	}
}
