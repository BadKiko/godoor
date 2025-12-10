package models

import "time"

// CookLog model matching Tandoor CookLog
type CookLog struct {
	BaseModel
	RecipeID   uint    `json:"-" gorm:"not null"`
	Recipe     Recipe  `json:"recipe" gorm:"foreignKey:RecipeID;references:ID"`
	Rating     *int    `json:"rating"`
	Servings   *int    `json:"servings"`
	Comment    *string `json:"comment"`

	CreatedByID uint `json:"-" gorm:"not null"`
	CreatedBy   User `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`
	SpaceID     uint `json:"-" gorm:"not null"`
	Space       Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
