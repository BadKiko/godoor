package serializers

import (
	"godoor/models"
	"time"
)

// StepFullSerializer matching Tandoor StepSerializer (full version for step endpoints)
type StepFullSerializer struct {
	ID                   uint                   `json:"id"`
	Name                 string                 `json:"name"`
	Instruction          string                 `json:"instruction"`
	Ingredients          []IngredientSerializer `json:"ingredients"`
	InstructionsMarkdown string                 `json:"instructions_markdown"`
	Time                 int                    `json:"time"`
	Order                int                    `json:"order"`
	ShowAsHeader         bool                   `json:"show_as_header"`
	ShowIngredientsTable bool                   `json:"show_ingredients_table"`
	File                 interface{}            `json:"file"`
	StepRecipe           interface{}            `json:"step_recipe"`
	StepRecipeData       interface{}            `json:"step_recipe_data"`
	NumRecipe            int                    `json:"numrecipe"`
}

// SerializeStep converts Step model to serializer
func SerializeStep(step *models.Step) StepFullSerializer {
	return StepFullSerializer{
		ID:                   step.ID,
		Name:                 step.Name,
		Instruction:          step.Instruction,
		Ingredients:          []IngredientSerializer{}, // TODO: implement ingredients
		InstructionsMarkdown: step.Instruction,          // TODO: implement markdown rendering
		Time:                 step.Time,
		Order:                step.Order,
		ShowAsHeader:         step.ShowAsHeader,
		ShowIngredientsTable: step.ShowIngredientsTable,
		File:                 nil, // TODO: implement file
		StepRecipe:           nil, // TODO: implement step recipe
		StepRecipeData:       nil, // TODO: implement step recipe data
		NumRecipe:            0,  // TODO: implement num recipe
	}
}

// StepListResponse represents paginated step response
type StepListResponse struct {
	Count     int                  `json:"count"`
	Next      *string              `json:"next"`
	Previous  *string              `json:"previous"`
	Results   []StepFullSerializer `json:"results"`
	Timestamp string               `json:"timestamp"`
}

// SerializeSteps converts slice of Step models to paginated response
func SerializeSteps(steps []models.Step, totalCount int) StepListResponse {
	result := make([]StepFullSerializer, len(steps))
	for i, step := range steps {
		result[i] = SerializeStep(&step)
	}
	return StepListResponse{
		Count:     totalCount,
		Next:      nil, // TODO: implement pagination
		Previous:  nil, // TODO: implement pagination
		Results:   result,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
