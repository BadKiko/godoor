package serializers

import (
	"godoor/config"
	"godoor/models"
	"time"

	"github.com/gomarkdown/markdown"
)

// getRecipeImageURL returns full image URL or nil if no image
func getRecipeImageURL(imagePath *string) interface{} {
	if imagePath != nil && *imagePath != "" {
		// If it's already a full URL (starts with http), return as is
		if len(*imagePath) > 4 && (*imagePath)[:4] == "http" {
			return *imagePath
		}
		// Otherwise, prepend MEDIA_URL
		return config.MediaURL + *imagePath
	}
	return nil
}

// IngredientSerializer matching Tandoor IngredientSerializer
type IngredientSerializer struct {
	ID                  uint            `json:"id"`
	Food                *FoodSerializer `json:"food"`
	Unit                *UnitSerializer `json:"unit"`
	Amount              float64         `json:"amount"`
	Conversions         []interface{}   `json:"conversions"` // TODO: implement conversions
	Note                string          `json:"note"`
	Order               int             `json:"order"`
	IsHeader            bool            `json:"is_header"`
	NoAmount            bool            `json:"no_amount"`
	OriginalText        string          `json:"original_text"`
	UsedInRecipes       []interface{}   `json:"used_in_recipes"` // TODO: implement
	AlwaysUsePluralUnit bool            `json:"always_use_plural_unit"`
	AlwaysUsePluralFood bool            `json:"always_use_plural_food"`
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
	ID                     uint                     `json:"id"`
	Name                   string                   `json:"name"`
	Description            string                   `json:"description,omitempty"`
	Image                  interface{}              `json:"image"`
	Keywords               []KeywordLabelSerializer `json:"keywords"`
	Steps                  []StepFullSerializer     `json:"steps"`
	WorkingTime            int                      `json:"working_time"`
	WaitingTime            int                      `json:"waiting_time"`
	CreatedBy              UserSerializer           `json:"created_by"`
	CreatedAt              time.Time                `json:"created_at"`
	UpdatedAt              time.Time                `json:"updated_at"`
	SourceURL              string                   `json:"source_url"`
	Internal               bool                     `json:"internal"`
	ShowIngredientOverview bool                     `json:"show_ingredient_overview"`
	Nutrition              interface{}              `json:"nutrition"` // TODO: implement NutritionInformation
	Properties             []PropertySerializer     `json:"properties"`
	FoodProperties         *map[string]interface{}  `json:"food_properties"` // TODO: implement food properties calculation
	Servings               int                      `json:"servings"`
	FilePath               string                   `json:"file_path"`
	ServingsText           string                   `json:"servings_text"`
	Rating                 *float64                 `json:"rating,omitempty"`
	LastCooked             *time.Time               `json:"last_cooked,omitempty"`
	Private                bool                     `json:"private"`
	Shared                 []interface{}            `json:"shared"` // TODO: implement shared users
}

// SerializeIngredient converts Ingredient model to serializer
func SerializeIngredient(ingredient *models.Ingredient) IngredientSerializer {
	return IngredientSerializer{
		ID:                  ingredient.ID,
		Food:                SerializeFood(ingredient.Food),
		Unit:                SerializeUnit(ingredient.Unit),
		Amount:              ingredient.Amount,
		Conversions:         []interface{}{}, // TODO: implement conversions
		Note:                ingredient.Note,
		Order:               ingredient.Order,
		IsHeader:            ingredient.IsHeader,
		NoAmount:            ingredient.NoAmount,
		OriginalText:        "",              // TODO: implement original text
		UsedInRecipes:       []interface{}{}, // TODO: implement used in recipes
		AlwaysUsePluralUnit: false,           // TODO: implement plural logic
		AlwaysUsePluralFood: false,           // TODO: implement plural logic
	}
}

// SerializeRecipe converts Recipe model to full serializer
func SerializeRecipe(recipe *models.Recipe) RecipeSerializer {
	description := ""
	if recipe.Description != nil {
		description = *recipe.Description
	}

	// Load keywords for recipe - they should be preloaded by the query
	var keywords []KeywordLabelSerializer
	for _, keyword := range recipe.Keywords {
		keywords = append(keywords, SerializeKeywordLabel(&keyword))
	}

	// Serialize steps if loaded
	steps := make([]StepFullSerializer, len(recipe.Steps))
	for i, step := range recipe.Steps {
		steps[i] = SerializeStep(&step)
	}

	// TODO: implement properties serialization
	properties := []PropertySerializer{}

	return RecipeSerializer{
		ID:                     recipe.ID,
		Name:                   recipe.Name,
		Description:            description,
		Image:                  getRecipeImageURL(recipe.Image),
		Keywords:               keywords,
		Steps:                  steps,
		WorkingTime:            recipe.WorkingTime,
		WaitingTime:            recipe.WaitingTime,
		CreatedBy:              SerializeUser(&recipe.CreatedBy),
		CreatedAt:              recipe.CreatedAt,
		UpdatedAt:              recipe.UpdatedAt,
		SourceURL:              "", // TODO: implement source URL
		Internal:               recipe.Internal,
		ShowIngredientOverview: true, // TODO: implement this field
		Nutrition:              nil,  // TODO: implement nutrition
		Properties:             properties,
		FoodProperties:         &map[string]interface{}{}, // TODO: implement food properties calculation
		Servings:               recipe.Servings,
		FilePath:               "", // TODO: implement file path
		ServingsText:           recipe.ServingsText,
		Rating:                 recipe.Rating,
		LastCooked:             nil, // TODO: implement last cooked
		Private:                recipe.Private,
		Shared:                 []interface{}{}, // TODO: implement shared users
	}
}

// SerializeRecipeWithSteps creates a recipe serializer with provided steps (for update responses)
func SerializeRecipeWithSteps(recipe *models.Recipe, requestSteps []map[string]interface{}) RecipeSerializer {
	base := SerializeRecipe(recipe)

	// Convert request steps to StepFullSerializer format (basic conversion)
	steps := make([]StepFullSerializer, len(requestSteps))
	for i, stepData := range requestSteps {
		instruction := getStringFromMap(stepData, "instruction")
		// Render markdown to HTML
		markdownBytes := []byte(instruction)
		htmlBytes := markdown.ToHTML(markdownBytes, nil, nil)
		instructionsHTML := string(htmlBytes)

		step := StepFullSerializer{
			Name:                 getStringFromMap(stepData, "name"),
			Instruction:          instruction,
			InstructionsMarkdown: instructionsHTML,
			Time:                 getIntFromMap(stepData, "time"),
			Order:                getIntFromMap(stepData, "order"),
			ShowAsHeader:         true, // default
			ShowIngredientsTable: true, // default
		}

		// Handle ingredients array - always initialize as empty slice if field exists
		if _, exists := stepData["ingredients"]; exists {
			if ingredientsData, ok := stepData["ingredients"]; ok {
				if ingredientsSlice, ok := ingredientsData.([]interface{}); ok {
					ingredients := make([]IngredientSerializer, len(ingredientsSlice))
					for j, ingData := range ingredientsSlice {
						if ingMap, ok := ingData.(map[string]interface{}); ok {
							ingredients[j] = IngredientSerializer{
								Amount:       getFloatFromMap(ingMap, "amount"),
								Note:         getStringFromMap(ingMap, "note"),
								Order:        getIntFromMap(ingMap, "order"),
								IsHeader:     getBoolFromMap(ingMap, "is_header"),
								NoAmount:     getBoolFromMap(ingMap, "no_amount"),
								OriginalText: getStringFromMap(ingMap, "original_text"),
								// TODO: implement Food, Unit, Conversions, UsedInRecipes
								Food:                nil,
								Unit:                nil,
								Conversions:         []interface{}{},
								UsedInRecipes:       []interface{}{},
								AlwaysUsePluralUnit: false,
								AlwaysUsePluralFood: false,
							}
						}
					}
					step.Ingredients = ingredients
				} else {
					// If ingredients field exists but is not an array, initialize as empty array
					step.Ingredients = []IngredientSerializer{}
				}
			} else {
				// If ingredients field exists, initialize as empty array
				step.Ingredients = []IngredientSerializer{}
			}
		}

		steps[i] = step
	}

	base.Steps = steps
	return base
}

// Helper functions to extract values from map
func getStringFromMap(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getIntFromMap(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
	}
	return 0
}

func getFloatFromMap(data map[string]interface{}, key string) float64 {
	if val, ok := data[key]; ok {
		if num, ok := val.(float64); ok {
			return num
		}
	}
	return 0.0
}

func getBoolFromMap(data map[string]interface{}, key string) bool {
	if val, ok := data[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
