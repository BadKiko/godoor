package models

import (
	"time"
)

// MealType model matching Tandoor MealType
type MealType struct {
	BaseModel
	Name      string    `json:"name" gorm:"not null"`
	Order     int       `json:"order" gorm:"default:0"`
	Color     *string   `json:"color"`
	Time      *string   `json:"time"` // Time field as string for simplicity
	Default   bool      `json:"default" gorm:"default:false"`
	CreatedByID uint    `json:"-" gorm:"not null"`
	CreatedBy User      `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`

	SpaceID uint  `json:"-" gorm:"not null"`
	Space   Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MealPlan model matching Tandoor MealPlan
type MealPlan struct {
	BaseModel
	RecipeID   *uint    `json:"-" gorm:"column:recipe_id"`
	Recipe     *Recipe  `json:"recipe,omitempty" gorm:"foreignKey:RecipeID;references:ID"`
	Servings   float64  `json:"servings" gorm:"default:1"`
	Title      string   `json:"title" gorm:"default:''"`
	CreatedByID uint    `json:"-" gorm:"not null"`
	CreatedBy  User     `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`
	MealTypeID uint     `json:"-" gorm:"not null"`
	MealType   MealType `json:"meal_type" gorm:"foreignKey:MealTypeID;references:ID"`
	Note       string   `json:"note" gorm:"default:''"`
	FromDate   time.Time `json:"from_date"`
	ToDate     time.Time `json:"to_date"`

	// TODO: Add shared ManyToMany relationship with User (related_name='plan_share')

	SpaceID uint  `json:"-" gorm:"not null"`
	Space   Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateMealType creates a new meal type
func CreateMealType(user *User, space *Space, name string, order int, color, mealTime *string, isDefault bool) (*MealType, error) {
	mealType := MealType{
		Name:        name,
		Order:       order,
		Color:       color,
		Time:        mealTime,
		Default:     isDefault,
		CreatedByID: user.ID,
		CreatedBy:   *user,
		SpaceID:     space.ID,
		Space:       *space,
	}

	if err := DB.Create(&mealType).Error; err != nil {
		return nil, err
	}

	return &mealType, nil
}

// CreateMealPlan creates a new meal plan
func CreateMealPlan(user *User, space *Space, mealType *MealType, fromDate, toDate time.Time, title, note string, recipeID *uint, servings float64) (*MealPlan, error) {
	mealPlan := MealPlan{
		RecipeID:    recipeID,
		Servings:    servings,
		Title:       title,
		CreatedByID: user.ID,
		CreatedBy:   *user,
		MealTypeID:  mealType.ID,
		MealType:    *mealType,
		Note:        note,
		FromDate:    fromDate,
		ToDate:      toDate,
		SpaceID:     space.ID,
		Space:       *space,
	}

	if err := DB.Create(&mealPlan).Error; err != nil {
		return nil, err
	}

	return &mealPlan, nil
}
