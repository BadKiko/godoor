package models

import (
	"time"
)

// ShoppingListRecipe model matching Tandoor ShoppingListRecipe
type ShoppingListRecipe struct {
	BaseModel
	Name     string  `json:"name" gorm:"default:''"`
	Servings float64 `json:"servings" gorm:"default:1"`

	RecipeID   *uint     `json:"-" gorm:"index"`
	Recipe     *Recipe   `json:"recipe,omitempty" gorm:"foreignKey:RecipeID;references:ID"`
	MealPlanID *uint     `json:"-" gorm:"index"`
	MealPlan   *MealPlan `json:"mealplan,omitempty" gorm:"foreignKey:MealPlanID;references:ID"`

	CreatedByID uint  `json:"-" gorm:"not null"`
	CreatedBy   User  `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`
	SpaceID     uint  `json:"-" gorm:"not null"`
	Space       Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ShoppingList model matching Tandoor ShoppingList
type ShoppingList struct {
	BaseModel
	Name        string  `json:"name" gorm:"default:''"`
	Description string  `json:"description" gorm:"type:text;default:''"`
	Color       *string `json:"color" gorm:"size:7"`

	SpaceID uint  `json:"-" gorm:"not null"`
	Space   Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ShoppingListEntry model matching Tandoor ShoppingListEntry
type ShoppingListEntry struct {
	BaseModel
	ShoppingLists []ShoppingList `json:"shopping_lists" gorm:"many2many:shopping_list_entry_lists;"`

	ListRecipeID *uint               `json:"-" gorm:"index"`
	ListRecipe   *ShoppingListRecipe `json:"list_recipe,omitempty" gorm:"foreignKey:ListRecipeID;references:ID"`

	FoodID       uint        `json:"-" gorm:"not null"`
	Food         Food        `json:"food" gorm:"foreignKey:FoodID;references:ID"`
	UnitID       *uint       `json:"-" gorm:"index"`
	Unit         *Unit       `json:"unit,omitempty" gorm:"foreignKey:UnitID;references:ID"`
	IngredientID *uint       `json:"-" gorm:"index"`
	Ingredient   *Ingredient `json:"ingredient,omitempty" gorm:"foreignKey:IngredientID;references:ID"`

	Amount      float64    `json:"amount" gorm:"default:0"`
	Order       int        `json:"order" gorm:"default:0"`
	Checked     bool       `json:"checked" gorm:"default:false"`
	CompletedAt *time.Time `json:"completed_at"`
	DelayUntil  *time.Time `json:"delay_until"`

	CreatedByID uint  `json:"-" gorm:"not null"`
	CreatedBy   User  `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`
	SpaceID     uint  `json:"-" gorm:"not null"`
	Space       Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
