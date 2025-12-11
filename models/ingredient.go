package models

import (
	"time"
)

// Ingredient model matching Tandoor Ingredient
type Ingredient struct {
	BaseModel
	FoodID   *uint   `json:"-" gorm:"index"`
	Food     *Food   `json:"food,omitempty" gorm:"foreignKey:FoodID;references:ID"`
	UnitID   *uint   `json:"-" gorm:"index"`
	Unit     *Unit   `json:"unit,omitempty" gorm:"foreignKey:UnitID;references:ID"`
	Amount   float64 `json:"amount" gorm:"type:decimal(16,16);default:0"`
	Note     string  `json:"note" gorm:"default:''"`
	IsHeader bool    `json:"is_header" gorm:"default:false"`
	NoAmount bool    `json:"no_amount" gorm:"default:false"`
	Order    int     `json:"order" gorm:"default:0"`

	SpaceID uint  `json:"-" gorm:"not null"`
	Space   Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	// Many-to-many with Steps through step_ingredients table
	Steps []Step `json:"-" gorm:"many2many:step_ingredients;"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateIngredient creates a new ingredient
func CreateIngredient(space *Space, food *Food, unit *Unit, amount float64, note string, order int) (*Ingredient, error) {
	ingredient := Ingredient{
		Food:    food,
		Unit:    unit,
		Amount:  amount,
		Note:    note,
		Order:   order,
		SpaceID: space.ID,
		Space:   *space,
	}

	if food != nil {
		ingredient.FoodID = &food.ID
	}
	if unit != nil {
		ingredient.UnitID = &unit.ID
	}

	if err := DB.Create(&ingredient).Error; err != nil {
		return nil, err
	}

	return &ingredient, nil
}
