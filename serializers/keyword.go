package serializers

import (
	"godoor/models"
)

// KeywordSerializer matching Tandoor KeywordSerializer
type KeywordSerializer struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SerializeKeyword converts Keyword model to serializer
func SerializeKeyword(keyword *models.Keyword) KeywordSerializer {
	return KeywordSerializer{
		ID:          keyword.ID,
		Name:        keyword.Name,
		Description: keyword.Description,
	}
}

// KeywordListResponse represents paginated keyword response matching Django REST Framework
type KeywordListResponse struct {
	Count    int                `json:"count"`
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
	Results  []KeywordSerializer `json:"results"`
}

// SerializeKeywords converts slice of Keyword models to paginated response
func SerializeKeywords(keywords []models.Keyword) KeywordListResponse {
	result := make([]KeywordSerializer, len(keywords))
	for i, keyword := range keywords {
		result[i] = SerializeKeyword(&keyword)
	}
	return KeywordListResponse{
		Count:    len(keywords),
		Next:     nil, // TODO: implement pagination
		Previous: nil, // TODO: implement pagination
		Results:  result,
	}
}
