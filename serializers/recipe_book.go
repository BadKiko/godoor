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

// SerializeRecipeBooks converts slice of RecipeBook models to slice of serializers
func SerializeRecipeBooks(books []models.RecipeBook) []RecipeBookSerializer {
	result := make([]RecipeBookSerializer, len(books))
	for i, book := range books {
		result[i] = SerializeRecipeBook(&book)
	}
	return result
}
