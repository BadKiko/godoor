package serializers

import (
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
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	Image       interface{} `json:"image"`
	Parent      *uint     `json:"parent"`
	NumChild    int       `json:"numchild"`
	NumRecipe   int       `json:"numrecipe"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	FullName    string    `json:"full_name"`
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
	Count     int                `json:"count"`
	Next      *string            `json:"next"`
	Previous  *string            `json:"previous"`
	Results   []KeywordSerializer `json:"results"`
	Timestamp string             `json:"timestamp"`
}

// SerializeKeywords converts slice of Keyword models to paginated response
func SerializeKeywords(keywords []models.Keyword) KeywordListResponse {
	result := make([]KeywordSerializer, len(keywords))
	for i, keyword := range keywords {
		result[i] = SerializeKeyword(&keyword)
	}
	return KeywordListResponse{
		Count:     len(keywords),
		Next:      nil, // TODO: implement pagination
		Previous:  nil, // TODO: implement pagination
		Results:   result,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
