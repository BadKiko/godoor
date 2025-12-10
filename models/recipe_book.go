package models

import (
	"time"
)

// RecipeBook model matching Tandoor RecipeBook
type RecipeBook struct {
	BaseModel
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description" gorm:"default:''"`
	CreatedByID uint   `json:"-" gorm:"not null"`
	CreatedBy   User   `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`
	Order       int    `json:"order" gorm:"default:0"`

	SpaceID uint  `json:"-" gorm:"not null"`
	Space   Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	// TODO: Add shared ManyToMany relationship with User
	// TODO: Add filter ForeignKey relationship with CustomFilter

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RecipeBookEntry represents the relationship between recipes and books
type RecipeBookEntry struct {
	BaseModel
	RecipeID uint   `json:"-" gorm:"not null"`
	Recipe   Recipe `json:"recipe,omitempty" gorm:"foreignKey:RecipeID;references:ID"`
	BookID   uint   `json:"-" gorm:"not null"`
	Book     RecipeBook `json:"book,omitempty" gorm:"foreignKey:BookID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateRecipeBook creates a new recipe book for a user
func CreateRecipeBook(user *User, space *Space, name, description string) (*RecipeBook, error) {
	book := RecipeBook{
		Name:        name,
		Description: description,
		CreatedByID: user.ID,
		CreatedBy:   *user,
		SpaceID:     space.ID,
		Space:       *space,
		Order:       0, // TODO: calculate order
	}

	if err := DB.Create(&book).Error; err != nil {
		return nil, err
	}

	return &book, nil
}
