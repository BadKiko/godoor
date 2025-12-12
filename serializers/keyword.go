package serializers

import (
	"fmt"
	"godoor/models"
	"time"
)

// KeywordLabelSerializer - simplified keyword for recipe keywords field
type KeywordLabelSerializer struct {
	ID    uint   `json:"id"`
	Label string `json:"label"`
}

// KeywordSerializer matching Tandoor KeywordSerializer
type KeywordSerializer struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Label       string      `json:"label"`
	Description string      `json:"description"`
	Image       interface{} `json:"image"`
	Parent      *uint       `json:"parent"`
	NumChild    int         `json:"numchild"`
	NumRecipe   int         `json:"numrecipe"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	FullName    string      `json:"full_name"`
}

// SerializeKeywordLabel converts Keyword to label serializer
func SerializeKeywordLabel(keyword *models.Keyword) KeywordLabelSerializer {
	return KeywordLabelSerializer{
		ID:    keyword.ID,
		Label: keyword.Name, // label is just the name
	}
}

// SerializeKeyword converts Keyword model to serializer
func SerializeKeyword(keyword *models.Keyword) KeywordSerializer {
	return KeywordSerializer{
		ID:          keyword.ID,
		Name:        keyword.Name,
		Label:       keyword.Name, // label is just the name
		Description: keyword.Description,
		Image:       nil, // TODO: implement image
		Parent:      nil, // TODO: implement tree parent
		NumChild:    0,   // TODO: implement child count
		NumRecipe:   0,   // TODO: implement recipe count
		CreatedAt:   keyword.CreatedAt,
		UpdatedAt:   keyword.UpdatedAt,
		FullName:    keyword.Name, // TODO: implement full tree name
	}
}

// KeywordListResponse represents paginated keyword response matching Django REST Framework
type KeywordListResponse struct {
	Count     int                 `json:"count"`
	Next      *string             `json:"next"`
	Previous  *string             `json:"previous"`
	Results   []KeywordSerializer `json:"results"`
	Timestamp string              `json:"timestamp"`
}

// SerializeKeywords converts slice of Keyword models to simple response (for backward compatibility)
func SerializeKeywords(keywords []models.Keyword) KeywordListResponse {
	result := make([]KeywordSerializer, len(keywords))
	for i, keyword := range keywords {
		result[i] = SerializeKeyword(&keyword)
	}
	return KeywordListResponse{
		Count:     len(keywords),
		Next:      nil,
		Previous:  nil,
		Results:   result,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// SerializeKeywordsPaginated converts slice of Keyword models to paginated response
func SerializeKeywordsPaginated(keywords []models.Keyword, totalCount int, page int, pageSize int) KeywordListResponse {
	result := make([]KeywordSerializer, len(keywords))
	for i, keyword := range keywords {
		result[i] = SerializeKeyword(&keyword)
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

	return KeywordListResponse{
		Count:     totalCount,
		Next:      next,
		Previous:  previous,
		Results:   result,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
