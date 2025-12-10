package serializers

import (
	"godoor/models"
)

// MealTypeSerializer matching Tandoor MealTypeSerializer
type MealTypeSerializer struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	Time      *string `json:"time"`
	Color     *string `json:"color"`
	Default   bool   `json:"default"`
	CreatedBy UserSerializer `json:"created_by"`
}

// SerializeMealType converts MealType model to serializer
func SerializeMealType(mealType *models.MealType) MealTypeSerializer {
	return MealTypeSerializer{
		ID:        mealType.ID,
		Name:      mealType.Name,
		Order:     mealType.Order,
		Time:      mealType.Time,
		Color:     mealType.Color,
		Default:   mealType.Default,
		CreatedBy: SerializeUser(&mealType.CreatedBy),
	}
}

// MealTypeListResponse represents paginated meal type response matching Django REST Framework
type MealTypeListResponse struct {
	Count    int                   `json:"count"`
	Next     *string               `json:"next"`
	Previous *string               `json:"previous"`
	Results  []MealTypeSerializer   `json:"results"`
}

// SerializeMealTypes converts slice of MealType models to paginated response
func SerializeMealTypes(mealTypes []models.MealType) MealTypeListResponse {
	result := make([]MealTypeSerializer, len(mealTypes))
	for i, mealType := range mealTypes {
		result[i] = SerializeMealType(&mealType)
	}
	return MealTypeListResponse{
		Count:    len(mealTypes),
		Next:     nil, // TODO: implement pagination
		Previous: nil, // TODO: implement pagination
		Results:  result,
	}
}

// RecipeOverviewSerializer moved to recipe.go

// MealPlanSerializer matching Tandoor MealPlanSerializer (simplified)
type MealPlanSerializer struct {
	ID             uint             `json:"id"`
	Title          string           `json:"title"`
	Recipe         interface{}      `json:"recipe,omitempty"`
	Servings       float64          `json:"servings"`
	Note           string           `json:"note"`
	MealType       MealTypeSerializer `json:"meal_type"`
	MealTypeName   string           `json:"meal_type_name"` // For compatibility
	FromDate       string           `json:"from_date"`
	ToDate         string           `json:"to_date"`
	Shared         []UserSerializer `json:"shared,omitempty"`
	Shopping       bool             `json:"shopping"` // TODO: implement shopping logic
}

// SerializeMealPlan converts MealPlan model to serializer
func SerializeMealPlan(mealPlan *models.MealPlan) MealPlanSerializer {
	var recipe interface{}
	if mealPlan.Recipe != nil {
		// Use the function from recipe.go
		r := SerializeRecipeOverview(mealPlan.Recipe)
		recipe = r
	}

	return MealPlanSerializer{
		ID:           mealPlan.ID,
		Title:        mealPlan.Title,
		Recipe:       recipe,
		Servings:     mealPlan.Servings,
		Note:         mealPlan.Note,
		MealType:     SerializeMealType(&mealPlan.MealType),
		MealTypeName: mealPlan.MealType.Name,
		FromDate:     mealPlan.FromDate.Format("2006-01-02T15:04:05Z07:00"),
		ToDate:       mealPlan.ToDate.Format("2006-01-02T15:04:05Z07:00"),
		Shared:       []UserSerializer{}, // TODO: implement shared users
		Shopping:     false,               // TODO: implement shopping logic
	}
}

// MealPlanListResponse represents paginated meal plan response matching Django REST Framework
type MealPlanListResponse struct {
	Count    int                   `json:"count"`
	Next     *string               `json:"next"`
	Previous *string               `json:"previous"`
	Results  []MealPlanSerializer   `json:"results"`
}

// SerializeMealPlans converts slice of MealPlan models to paginated response
func SerializeMealPlans(mealPlans []models.MealPlan) MealPlanListResponse {
	result := make([]MealPlanSerializer, len(mealPlans))
	for i, mealPlan := range mealPlans {
		result[i] = SerializeMealPlan(&mealPlan)
	}
	return MealPlanListResponse{
		Count:    len(mealPlans),
		Next:     nil, // TODO: implement pagination
		Previous: nil, // TODO: implement pagination
		Results:  result,
	}
}
