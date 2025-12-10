package serializers

import (
	"godoor/models"
	"time"
)

// CookLogSerializer matching Tandoor CookLogSerializer
type CookLogSerializer struct {
	ID        uint            `json:"id"`
	Recipe    RecipeSerializer `json:"recipe"`
	Servings  *int            `json:"servings"`
	Rating    *int            `json:"rating"`
	Comment   *string         `json:"comment"`
	CreatedBy UserSerializer  `json:"created_by"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// SerializeCookLog converts CookLog model to serializer
func SerializeCookLog(cookLog *models.CookLog) CookLogSerializer {
	comment := ""
	if cookLog.Comment != nil {
		comment = *cookLog.Comment
	}

	return CookLogSerializer{
		ID:        cookLog.ID,
		Recipe:    SerializeRecipe(&cookLog.Recipe),
		Servings:  cookLog.Servings,
		Rating:    cookLog.Rating,
		Comment:   &comment,
		CreatedBy: SerializeUser(&cookLog.CreatedBy),
		CreatedAt: cookLog.CreatedAt,
		UpdatedAt: cookLog.UpdatedAt,
	}
}

// CookLogListResponse represents paginated cook log response
type CookLogListResponse struct {
	Count     int                `json:"count"`
	Next      *string            `json:"next"`
	Previous  *string            `json:"previous"`
	Results   []CookLogSerializer `json:"results"`
	Timestamp string             `json:"timestamp"`
}

// SerializeCookLogs converts slice of CookLog models to paginated response
func SerializeCookLogs(cookLogs []models.CookLog, totalCount int) CookLogListResponse {
	result := make([]CookLogSerializer, len(cookLogs))
	for i, cookLog := range cookLogs {
		result[i] = SerializeCookLog(&cookLog)
	}
	return CookLogListResponse{
		Count:     totalCount,
		Next:      nil, // TODO: implement pagination
		Previous:  nil, // TODO: implement pagination
		Results:   result,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
