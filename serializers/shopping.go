package serializers

import (
	"godoor/models"
	"time"
)

// ShoppingListRecipeSerializer matching Tandoor ShoppingListRecipeSerializer
type ShoppingListRecipeSerializer struct {
	ID       uint                      `json:"id"`
	Name     string                    `json:"name"`
	Servings float64                   `json:"servings"`
	Recipe   *RecipeOverviewSerializer `json:"recipe,omitempty"`
	MealPlan interface{}               `json:"mealplan,omitempty"` // TODO: implement meal plan serializer
}

// SerializeShoppingListRecipe converts ShoppingListRecipe model to serializer
func SerializeShoppingListRecipe(recipe *models.ShoppingListRecipe) ShoppingListRecipeSerializer {
	var recipeData *RecipeOverviewSerializer
	if recipe.Recipe != nil {
		r := SerializeRecipeOverview(recipe.Recipe)
		recipeData = &r
	}

	return ShoppingListRecipeSerializer{
		ID:       recipe.ID,
		Name:     recipe.Name,
		Servings: recipe.Servings,
		Recipe:   recipeData,
		MealPlan: nil, // TODO: implement
	}
}

// ShoppingListSerializer matching Tandoor ShoppingListSerializer
type ShoppingListSerializer struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       *string   `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SerializeShoppingList converts ShoppingList model to serializer
func SerializeShoppingList(list *models.ShoppingList) ShoppingListSerializer {
	return ShoppingListSerializer{
		ID:          list.ID,
		Name:        list.Name,
		Description: list.Description,
		Color:       list.Color,
		CreatedAt:   list.CreatedAt,
		UpdatedAt:   list.UpdatedAt,
	}
}

// ShoppingListEntrySerializer matching Tandoor ShoppingListEntrySerializer
type ShoppingListEntrySerializer struct {
	ID             uint                          `json:"id"`
	ListRecipe     *uint                         `json:"list_recipe"`
	ShoppingLists  []ShoppingListSerializer      `json:"shopping_lists"`
	Food           *FoodShoppingSerializer       `json:"food"`
	Unit           *UnitSerializer               `json:"unit,omitempty"`
	Amount         float64                       `json:"amount"`
	Order          int                           `json:"order"`
	Checked        bool                          `json:"checked"`
	Ingredient     *uint                         `json:"ingredient,omitempty"`
	ListRecipeData *ShoppingListRecipeSerializer `json:"list_recipe_data,omitempty"`
	CreatedBy      UserSerializer                `json:"created_by"`
	CreatedAt      time.Time                     `json:"created_at"`
	UpdatedAt      time.Time                     `json:"updated_at"`
	CompletedAt    *time.Time                    `json:"completed_at,omitempty"`
	DelayUntil     *time.Time                    `json:"delay_until,omitempty"`
}

// SerializeShoppingListEntry converts ShoppingListEntry model to serializer
func SerializeShoppingListEntry(entry *models.ShoppingListEntry) ShoppingListEntrySerializer {
	// Serialize shopping lists
	shoppingLists := make([]ShoppingListSerializer, len(entry.ShoppingLists))
	for i, list := range entry.ShoppingLists {
		shoppingLists[i] = SerializeShoppingList(&list)
	}

	// Serialize food
	var food *FoodShoppingSerializer
	if entry.FoodID > 0 {
		// If Food is not loaded, try to load it
		if entry.Food.ID == 0 {
			var foodModel models.Food
			if err := models.DB.First(&foodModel, entry.FoodID).Error; err == nil {
				entry.Food = foodModel
			}
		}
		if entry.Food.ID > 0 {
			f := SerializeFoodShopping(&entry.Food)
			food = &f
		}
	}

	// Serialize unit
	var unit *UnitSerializer
	if entry.Unit != nil {
		unit = SerializeUnit(entry.Unit)
	}

	// Serialize list recipe data
	var listRecipeData *ShoppingListRecipeSerializer
	if entry.ListRecipe != nil {
		lrd := SerializeShoppingListRecipe(entry.ListRecipe)
		listRecipeData = &lrd
	}

	// Handle list_recipe field
	var listRecipe *uint
	if entry.ListRecipeID != nil {
		listRecipe = entry.ListRecipeID
	}

	// Handle ingredient field
	var ingredient *uint
	if entry.IngredientID != nil {
		ingredient = entry.IngredientID
	}

	return ShoppingListEntrySerializer{
		ID:             entry.ID,
		ListRecipe:     listRecipe,
		ShoppingLists:  shoppingLists,
		Food:           food,
		Unit:           unit,
		Amount:         entry.Amount,
		Order:          entry.Order,
		Checked:        entry.Checked,
		Ingredient:     ingredient,
		ListRecipeData: listRecipeData,
		CreatedBy:      SerializeUser(&entry.CreatedBy),
		CreatedAt:      entry.CreatedAt,
		UpdatedAt:      entry.UpdatedAt,
		CompletedAt:    entry.CompletedAt,
		DelayUntil:     entry.DelayUntil,
	}
}

// ShoppingListEntryListResponse represents paginated shopping list entry response
type ShoppingListEntryListResponse struct {
	Count    int                           `json:"count"`
	Next     *string                       `json:"next"`
	Previous *string                       `json:"previous"`
	Results  []ShoppingListEntrySerializer `json:"results"`
}

// SerializeShoppingListEntries converts slice of ShoppingListEntry models to paginated response
func SerializeShoppingListEntries(entries []models.ShoppingListEntry, totalCount int, page int, pageSize int) ShoppingListEntryListResponse {
	result := make([]ShoppingListEntrySerializer, len(entries))
	for i, entry := range entries {
		result[i] = SerializeShoppingListEntry(&entry)
	}

	return ShoppingListEntryListResponse{
		Count:   totalCount,
		Results: result,
		// TODO: implement pagination URLs
	}
}
