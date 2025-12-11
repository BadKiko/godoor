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
}
