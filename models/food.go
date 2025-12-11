package models

import (
	"time"
)

// Food model matching Tandoor Food (simplified)
type Food struct {
	BaseModel
	Name        string  `json:"name" gorm:"not null"`
	PluralName  *string `json:"plural_name"`
	Description string  `json:"description" gorm:"default:''"`

	SpaceID uint  `json:"-" gorm:"not null"`
	Space   Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	// Many-to-many with shopping lists
	ShoppingLists []ShoppingList `json:"-" gorm:"many2many:shopping_list_entry_lists;"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateFood creates a new food item
func CreateFood(space *Space, name, description string) (*Food, error) {
	food := Food{
		Name:        name,
		Description: description,
		SpaceID:     space.ID,
		Space:       *space,
	}

	if err := DB.Create(&food).Error; err != nil {
		return nil, err
	}

	return &food, nil
}
