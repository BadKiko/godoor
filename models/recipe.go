package models

import (
	"time"
)

// Recipe model matching Tandoor Recipe (simplified)
type Recipe struct {
	BaseModel
	Name        string  `json:"name" gorm:"not null"`
	Description *string `json:"description"`
	Servings    int     `json:"servings" gorm:"default:1"`
	ServingsText string `json:"servings_text" gorm:"default:''"`
	WorkingTime int    `json:"working_time" gorm:"default:0"`
	WaitingTime int    `json:"waiting_time" gorm:"default:0"`
	Internal    bool   `json:"internal" gorm:"default:false"`
	Private     bool   `json:"private" gorm:"default:false"`

	CreatedByID uint   `json:"-" gorm:"not null"`
	CreatedBy   User   `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`
	SpaceID     uint   `json:"-" gorm:"not null"`
	Space       Space  `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	// One-to-many relationships
	Steps      []Step     `json:"steps" gorm:"foreignKey:RecipeID"`
	Properties []Property `json:"properties" gorm:"many2many:recipe_properties;"`

	// Rating field for sorting (simplified)
	Rating *float64 `json:"rating"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateRecipe creates a new recipe
func CreateRecipe(user *User, space *Space, name string, description *string, servings int) (*Recipe, error) {
	recipe := Recipe{
		Name:        name,
		Description: description,
		Servings:    servings,
		CreatedByID: user.ID,
		CreatedBy:   *user,
		SpaceID:     space.ID,
		Space:       *space,
	}

	if err := DB.Create(&recipe).Error; err != nil {
		return nil, err
	}

	return &recipe, nil
}
