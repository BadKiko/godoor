package models

import (
	"time"
)

// Unit model matching Tandoor Unit (simplified)
type Unit struct {
	BaseModel
	Name        string  `json:"name" gorm:"not null;uniqueIndex:idx_unit_space_name"`
	PluralName  *string `json:"plural_name"`
	Description string  `json:"description" gorm:"default:''"`

	SpaceID uint  `json:"-" gorm:"not null"`
	Space   Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUnit creates a new unit
func CreateUnit(space *Space, name, description string) (*Unit, error) {
	unit := Unit{
		Name:        name,
		Description: description,
		SpaceID:     space.ID,
		Space:       *space,
	}

	if err := DB.Create(&unit).Error; err != nil {
		return nil, err
	}

	return &unit, nil
}
