package models

import "time"

// PropertyType model matching Tandoor PropertyType
type PropertyType struct {
	BaseModel
	Name        string `json:"name" gorm:"not null"`
	Unit        string `json:"unit" gorm:"not null"`
	Description string `json:"description" gorm:"type:text"`

	// Property type choices
	PropertyType string `json:"property_type" gorm:"default:'OTHER'"` // NUTRITION, ALLERGEN, PRICE, GOAL, OTHER

	SpaceID uint `json:"-" gorm:"not null"`
	Space   Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Property model matching Tandoor Property
type Property struct {
	BaseModel
	PropertyAmount    *float64     `json:"property_amount"`
	PropertyTypeID    uint         `json:"-" gorm:"not null"`
	PropertyType      PropertyType `json:"property_type" gorm:"foreignKey:PropertyTypeID;references:ID"`
	OpenDataFoodSlug  *string      `json:"open_data_food_slug"`

	SpaceID uint `json:"-" gorm:"not null"`
	Space   Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	// Many-to-many with recipes
	Recipes []Recipe `json:"-" gorm:"many2many:recipe_properties;"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
