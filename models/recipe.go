package models

import (
	"time"
)

// Recipe model - placeholder for now
type Recipe struct {
	BaseModel
	Name   string `json:"name" gorm:"not null"`
	SpaceID uint  `json:"-" gorm:"not null"`
	Space   Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
