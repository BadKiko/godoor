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
}
