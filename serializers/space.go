package serializers

import (
	"godoor/models"
	"time"
)

// SpaceSerializer matching Tandoor SpaceSerializer
type SpaceSerializer struct {
	ID                   uint          `json:"id"`
	Name                 string        `json:"name"`
	CreatedBy            *UserSerializer `json:"created_by,omitempty"`
	CreatedAt            time.Time     `json:"created_at"`
	UserCount            int           `json:"user_count"`
	RecipeCount          int           `json:"recipe_count"`
	FileSizeMb           float64       `json:"file_size_mb"`
	AiMonthlyCreditsUsed int           `json:"ai_monthly_credits_used"`
	Message              string        `json:"message"`
}

// SerializeSpace converts Space model to SpaceSerializer
func SerializeSpace(space *models.Space) SpaceSerializer {
	return SpaceSerializer{
		ID:                   space.ID,
		Name:                 space.Name,
		CreatedBy:            nil, // Will be set if needed
		CreatedAt:            space.CreatedAt,
		UserCount:            0,   // TODO: implement counting
		RecipeCount:          0,   // TODO: implement counting
		FileSizeMb:           0,   // TODO: implement file size calculation
		AiMonthlyCreditsUsed: 0,   // TODO: implement AI credits tracking
		Message:              space.Message,
	}
}

// SerializeSpaces converts slice of Space models to slice of SpaceSerializers
func SerializeSpaces(spaces []models.Space) []SpaceSerializer {
	result := make([]SpaceSerializer, len(spaces))
	for i, space := range spaces {
		result[i] = SerializeSpace(&space)
	}
	return result
}
