package models

import "time"

// Step model matching Tandoor Step
type Step struct {
	BaseModel
	Name                string `json:"name" gorm:"default:''"`
	Instruction         string `json:"instruction" gorm:"type:text"`
	Time                int    `json:"time" gorm:"default:0"`
	Order               int    `json:"order" gorm:"default:0"`
	ShowAsHeader        bool   `json:"show_as_header" gorm:"default:true"`
	ShowIngredientsTable bool  `json:"show_ingredients_table" gorm:"default:true"`

	SpaceID uint `json:"-" gorm:"not null"`
	Space   Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	// Many-to-many with recipes (for step recipes)
	Recipes []Recipe `json:"-" gorm:"many2many:recipe_steps;"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
