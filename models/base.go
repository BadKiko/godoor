package models

import (
	"time"

	"gorm.io/gorm"
)

var DB *gorm.DB

// Base model with common fields
type BaseModel struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AutoMigrate runs all migrations
func AutoMigrate() {
	err := DB.AutoMigrate(
		&User{},
		&Space{},
		&UserSpace{},
		&AccessToken{},
		&UserPreference{},
		&Recipe{},
		&RecipeBook{},
		&RecipeBookEntry{},
		&Keyword{},
		&MealType{},
		&MealPlan{},
		&Food{},
		&Unit{},
		&Ingredient{},
		&Step{},
		&Property{},
		&PropertyType{},
		&CookLog{},
		&ShoppingListRecipe{},
		&ShoppingList{},
		&ShoppingListEntry{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Create many-to-many junction tables that GORM might not create automatically
	err = DB.Exec(`CREATE TABLE IF NOT EXISTS recipe_keywords (
		recipe_id INTEGER,
		keyword_id INTEGER,
		PRIMARY KEY (recipe_id, keyword_id),
		FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
		FOREIGN KEY (keyword_id) REFERENCES keywords(id) ON DELETE CASCADE
	)`).Error
	if err != nil {
		panic("Failed to create recipe_keywords junction table: " + err.Error())
	}

	// Create performance indexes
	indexes := []string{
		// Recipes indexes
		"CREATE INDEX IF NOT EXISTS idx_recipes_space_id ON recipes(space_id)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_created_by_id ON recipes(created_by_id)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_space_created_at ON recipes(space_id, created_at)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_space_name ON recipes(space_id, name)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_space_rating ON recipes(space_id, rating)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_name ON recipes(name)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_description ON recipes(description)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_space_updated_at ON recipes(space_id, updated_at)",

		// Keywords indexes
		"CREATE INDEX IF NOT EXISTS idx_keywords_space_id ON keywords(space_id)",
		"CREATE INDEX IF NOT EXISTS idx_keywords_space_name ON keywords(space_id, name)",
		"CREATE INDEX IF NOT EXISTS idx_keywords_id ON keywords(id)",
		"CREATE INDEX IF NOT EXISTS idx_keywords_name ON keywords(name)",

		// Recipe-Keywords junction table indexes (critical for filtering)
		"CREATE INDEX IF NOT EXISTS idx_recipe_keywords_recipe_id ON recipe_keywords(recipe_id)",
		"CREATE INDEX IF NOT EXISTS idx_recipe_keywords_keyword_id ON recipe_keywords(keyword_id)",
		"CREATE INDEX IF NOT EXISTS idx_recipe_keywords_both ON recipe_keywords(recipe_id, keyword_id)",

		// Steps indexes
		"CREATE INDEX IF NOT EXISTS idx_steps_recipe_id ON steps(recipe_id)",
		"CREATE INDEX IF NOT EXISTS idx_steps_recipe_order ON steps(recipe_id, \"order\")",

		// Ingredients indexes
		"CREATE INDEX IF NOT EXISTS idx_ingredients_food_id ON ingredients(food_id)",
		"CREATE INDEX IF NOT EXISTS idx_ingredients_unit_id ON ingredients(unit_id)",
		"CREATE INDEX IF NOT EXISTS idx_ingredients_space_id ON ingredients(space_id)",

		// Step ingredients junction table indexes
		"CREATE INDEX IF NOT EXISTS idx_step_ingredients_step_id ON step_ingredients(step_id)",
		"CREATE INDEX IF NOT EXISTS idx_step_ingredients_ingredient_id ON step_ingredients(ingredient_id)",

		// Foods indexes
		"CREATE INDEX IF NOT EXISTS idx_foods_space_id ON foods(space_id)",
		"CREATE INDEX IF NOT EXISTS idx_foods_space_name ON foods(space_id, name)",

		// Units indexes
		"CREATE INDEX IF NOT EXISTS idx_units_space_id ON units(space_id)",
		"CREATE INDEX IF NOT EXISTS idx_units_space_name ON units(space_id, name)",

		// Cook logs indexes (critical for favorite count performance)
		"CREATE INDEX IF NOT EXISTS idx_cook_logs_space_id ON cook_logs(space_id)",
		"CREATE INDEX IF NOT EXISTS idx_cook_logs_recipe_id ON cook_logs(recipe_id)",
		"CREATE INDEX IF NOT EXISTS idx_cook_logs_space_recipe ON cook_logs(space_id, recipe_id)",
		"CREATE INDEX IF NOT EXISTS idx_cook_logs_created_by_id ON cook_logs(created_by_id)",
		"CREATE INDEX IF NOT EXISTS idx_cook_logs_recipe_created_by ON cook_logs(recipe_id, created_by_id)",
		"CREATE INDEX IF NOT EXISTS idx_cook_logs_space_created_by ON cook_logs(space_id, created_by_id)",

		// Meal plans indexes
		"CREATE INDEX IF NOT EXISTS idx_meal_plans_space_id ON meal_plans(space_id)",
		"CREATE INDEX IF NOT EXISTS idx_meal_plans_recipe_id ON meal_plans(recipe_id)",
		"CREATE INDEX IF NOT EXISTS idx_meal_plans_meal_type_id ON meal_plans(meal_type_id)",
		"CREATE INDEX IF NOT EXISTS idx_meal_plans_date_range ON meal_plans(from_date, to_date)",

		// Recipe book entries indexes
		"CREATE INDEX IF NOT EXISTS idx_recipe_book_entries_book_id ON recipe_book_entries(book_id)",
		"CREATE INDEX IF NOT EXISTS idx_recipe_book_entries_recipe_id ON recipe_book_entries(recipe_id)",

		// Shopping list indexes
		"CREATE INDEX IF NOT EXISTS idx_shopping_list_recipes_space_id ON shopping_list_recipes(space_id)",
		"CREATE INDEX IF NOT EXISTS idx_shopping_list_entries_list_recipe_id ON shopping_list_entries(list_recipe_id)",
		"CREATE INDEX IF NOT EXISTS idx_shopping_list_entries_ingredient_id ON shopping_list_entries(ingredient_id)",
		"CREATE INDEX IF NOT EXISTS idx_shopping_list_entries_unit_id ON shopping_list_entries(unit_id)",

		// Access tokens indexes (for faster auth)
		"CREATE INDEX IF NOT EXISTS idx_access_tokens_user_id ON access_tokens(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_access_tokens_expires ON access_tokens(expires)",

		// Spaces indexes
		"CREATE INDEX IF NOT EXISTS idx_spaces_created_by_id ON spaces(created_by_id)",
		"CREATE INDEX IF NOT EXISTS idx_spaces_id ON spaces(id)",

		// User spaces indexes (critical for space resolution)
		"CREATE INDEX IF NOT EXISTS idx_user_spaces_user_id ON user_spaces(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_user_spaces_space_id ON user_spaces(space_id)",
		"CREATE INDEX IF NOT EXISTS idx_user_spaces_user_active ON user_spaces(user_id, active)",
		"CREATE INDEX IF NOT EXISTS idx_user_spaces_active ON user_spaces(active)",

		// Users indexes
		"CREATE INDEX IF NOT EXISTS idx_users_id ON users(id)",
		"CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)",
	}

	for _, indexSQL := range indexes {
		if err := DB.Exec(indexSQL).Error; err != nil {
			// Log warning but don't fail migration
			println("Warning: Failed to create index:", indexSQL, "Error:", err.Error())
		}
	}

	// Analyze database for query optimization
	DB.Exec("ANALYZE")
}
