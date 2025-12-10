package models

import (
	"time"
)

// Keyword model matching Tandoor Keyword (simplified, no tree structure for now)
type Keyword struct {
	BaseModel
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description" gorm:"default:''"`

	SpaceID uint  `json:"-" gorm:"not null"`
	Space   Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateKeyword creates a new keyword
func CreateKeyword(space *Space, name, description string) (*Keyword, error) {
	keyword := Keyword{
		Name:        name,
		Description: description,
		SpaceID:     space.ID,
		Space:       *space,
	}

	if err := DB.Create(&keyword).Error; err != nil {
		return nil, err
	}

	return &keyword, nil
}
