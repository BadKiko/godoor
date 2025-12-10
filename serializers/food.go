package serializers

import (
	"godoor/models"
	"time"
)

// FoodSerializer matching Tandoor FoodSerializer (simplified)
type FoodSerializer struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       interface{} `json:"image"`
	Parent      *uint  `json:"parent"`
	NumChild    int    `json:"numchild"`
	NumRecipe   int    `json:"numrecipe"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	FullName    string `json:"full_name"`
}

// FoodListResponse represents paginated food response
type FoodListResponse struct {
	Count     int             `json:"count"`
	Next      *string         `json:"next"`
	Previous  *string         `json:"previous"`
	Results   []FoodSerializer `json:"results"`
	Timestamp string          `json:"timestamp"`
}

// SerializeFood converts Food model to serializer
func SerializeFood(food *models.Food) FoodSerializer {
	return FoodSerializer{
		ID:          food.ID,
		Name:        food.Name,
		Description: food.Description,
		Image:       nil, // TODO: implement image
		Parent:      nil, // TODO: implement tree parent
		NumChild:    0,   // TODO: implement child count
		NumRecipe:   0,   // TODO: implement recipe count
		CreatedAt:   food.CreatedAt,
		UpdatedAt:   food.UpdatedAt,
		FullName:    food.Name, // TODO: implement full tree name
	}
}

// SerializeFoods converts slice of Food models to paginated response
func SerializeFoods(foods []models.Food, totalCount int) FoodListResponse {
	result := make([]FoodSerializer, len(foods))
	for i, food := range foods {
		result[i] = SerializeFood(&food)
	}
	return FoodListResponse{
		Count:     totalCount,
		Next:      nil, // TODO: implement pagination
		Previous:  nil, // TODO: implement pagination
		Results:   result,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
