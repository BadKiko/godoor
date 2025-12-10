package serializers

import (
	"godoor/models"
	"time"
)

// IngredientSerializer placeholder (simplified)
type IngredientSerializer struct {
	ID          uint        `json:"id"`
	Amount      float64     `json:"amount"`
	Unit        interface{} `json:"unit"` // TODO: implement Unit
	Food        interface{} `json:"food"` // TODO: implement Food
	Note        string      `json:"note"`
	IsHeader    bool        `json:"is_header"`
	NoAmount    bool        `json:"no_amount"`
	Order       int         `json:"order"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// StepSerializer matching Tandoor StepSerializer
type StepSerializer struct {
	ID                   uint                   `json:"id"`
	Name                 string                 `json:"name"`
	Instruction          string                 `json:"instruction"`
	Ingredients          []IngredientSerializer `json:"ingredients"`
	InstructionsMarkdown string                 `json:"instructions_markdown"`
	Time                 int                    `json:"time"`
	Order                int                    `json:"order"`
	ShowAsHeader         bool                   `json:"show_as_header"`
	ShowIngredientsTable bool                   `json:"show_ingredients_table"`
	File                 interface{}            `json:"file"` // TODO: implement UserFile
	StepRecipe           interface{}            `json:"step_recipe"`
	StepRecipeData       interface{}            `json:"step_recipe_data"`
	NumRecipe            int                    `json:"numrecipe"`
}

// PropertyTypeSerializer matching Tandoor PropertyTypeSerializer
type PropertyTypeSerializer struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}

// PropertySerializer matching Tandoor PropertySerializer
type PropertySerializer struct {
	ID             uint                   `json:"id"`
	PropertyAmount *float64               `json:"property_amount"`
	PropertyType   PropertyTypeSerializer `json:"property_type"`
}

// RecipeSerializer full serializer matching Tandoor RecipeSerializer
type RecipeSerializer struct {
	ID                    uint              `json:"id"`
	Name                  string            `json:"name"`
	Description           string            `json:"description,omitempty"`
	Image                 interface{}       `json:"image"`
	Keywords              []KeywordLabelSerializer `json:"keywords"`
	Steps                 []StepSerializer `json:"steps"`
	WorkingTime           int               `json:"working_time"`
	WaitingTime           int               `json:"waiting_time"`
	CreatedBy             UserSerializer    `json:"created_by"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
	SourceURL             string            `json:"source_url"`
	Internal              bool              `json:"internal"`
	ShowIngredientOverview bool             `json:"show_ingredient_overview"`
	Nutrition             interface{}       `json:"nutrition"` // TODO: implement NutritionInformation
	Properties            []PropertySerializer `json:"properties"`
	FoodProperties        *map[string]interface{} `json:"food_properties"` // TODO: implement food properties calculation
	Servings              int               `json:"servings"`
	FilePath              string            `json:"file_path"`
	ServingsText          string            `json:"servings_text"`
	Rating                *float64          `json:"rating,omitempty"`
	LastCooked            *time.Time        `json:"last_cooked,omitempty"`
	Private               bool              `json:"private"`
	Shared                []interface{}     `json:"shared"` // TODO: implement shared users
}

// SerializeRecipe converts Recipe model to full serializer
func SerializeRecipe(recipe *models.Recipe) RecipeSerializer {
	description := ""
	if recipe.Description != nil {
		description = *recipe.Description
	}

	// TODO: implement proper keywords loading
	keywords := []KeywordLabelSerializer{}

	// TODO: implement steps serialization
	steps := []StepSerializer{}

	// TODO: implement properties serialization
	properties := []PropertySerializer{}

	return RecipeSerializer{
		ID:                    recipe.ID,
		Name:                  recipe.Name,
		Description:           description,
		Image:                 nil, // TODO: implement image
		Keywords:              keywords,
		Steps:                 steps,
		WorkingTime:           recipe.WorkingTime,
		WaitingTime:           recipe.WaitingTime,
		CreatedBy:             SerializeUser(&recipe.CreatedBy),
		CreatedAt:             recipe.CreatedAt,
		UpdatedAt:             recipe.UpdatedAt,
		SourceURL:             "", // TODO: implement source URL
		Internal:              recipe.Internal,
		ShowIngredientOverview: true, // TODO: implement this field
		Nutrition:             nil, // TODO: implement nutrition
		Properties:            properties,
		FoodProperties:        &map[string]interface{}{}, // TODO: implement food properties calculation
		Servings:              recipe.Servings,
		FilePath:              "", // TODO: implement file path
		ServingsText:          recipe.ServingsText,
		Rating:                recipe.Rating,
		LastCooked:            nil, // TODO: implement last cooked
		Private:               recipe.Private,
		Shared:                []interface{}{}, // TODO: implement shared users
	}
}
