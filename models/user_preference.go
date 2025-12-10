package models

import (
	"time"
)

// UserPreference model matching Tandoor UserPreference
type UserPreference struct {
	BaseModel
	UserID      uint   `json:"-" gorm:"not null"`
	User        User   `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	SpaceID     *uint  `json:"-" gorm:"column:space_id"`
	Space       *Space `json:"-" gorm:"foreignKey:SpaceID;references:ID"`

	// Theme settings
	Theme      string `json:"theme" gorm:"default:'TANDOOR'"`
	NavBgColor string `json:"nav_bg_color" gorm:"default:'#ddbf86'"`
	NavTextColor string `json:"nav_text_color" gorm:"default:'DARK'"`
	NavShowLogo bool   `json:"nav_show_logo" gorm:"default:true"`

	// Navigation settings
	NavSticky bool `json:"nav_sticky" gorm:"default:true"`

	// Space settings
	MaxOwnedSpaces int `json:"max_owned_spaces" gorm:"default:3"`

	// Units and display
	DefaultUnit        string  `json:"default_unit" gorm:"default:'g'"`
	UseFractions       bool    `json:"use_fractions" gorm:"default:true"`
	UseKj             bool    `json:"use_kj" gorm:"default:false"`
	IngredientDecimals int     `json:"ingredient_decimals" gorm:"default:2"`
	DefaultDelay       float64 `json:"default_delay" gorm:"default:4"`

	// Default page
	DefaultPage string `json:"default_page" gorm:"default:'SEARCH'"`

	// Comments and planning
	Comments                   bool `json:"comments" gorm:"default:true"`
	MealplanAutoaddShopping    bool `json:"mealplan_autoadd_shopping" gorm:"default:false"`
	MealplanAutoexcludeOnhand  bool `json:"mealplan_autoexclude_onhand" gorm:"default:true"`
	MealplanAutoincludeRelated bool `json:"mealplan_autoinclude_related" gorm:"default:true"`

	// Shopping settings
	ShoppingAutoSync      int  `json:"shopping_auto_sync" gorm:"default:5"`
	ShoppingRecentDays    int  `json:"shopping_recent_days" gorm:"default:7"`
	ShoppingAddOnhand     bool `json:"shopping_add_onhand" gorm:"default:false"`
	ShoppingUpdateFoodLists bool `json:"shopping_update_food_lists" gorm:"default:true"`

	// Supermarket settings
	FilterToSupermarket bool `json:"filter_to_supermarket" gorm:"default:false"`

	// UI settings
	LeftHanded         bool `json:"left_handed" gorm:"default:false"`
	ShowStepIngredients bool `json:"show_step_ingredients" gorm:"default:true"`

	// CSV settings
	CsvDelim  string `json:"csv_delim" gorm:"default:','"`
	CsvPrefix string `json:"csv_prefix" gorm:"default:''"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserPreference creates default user preferences for a user
func CreateUserPreference(user *User) (*UserPreference, error) {
	preference := UserPreference{
		UserID: user.ID,
		User:   *user,
		// All other fields will use their default values
	}

	if err := DB.Create(&preference).Error; err != nil {
		return nil, err
	}

	return &preference, nil
}

// GetOrCreateUserPreference gets existing preferences or creates default ones
func GetOrCreateUserPreference(user *User) (*UserPreference, error) {
	var preference UserPreference
	err := DB.Where("user_id = ?", user.ID).Preload("User").First(&preference).Error
	if err != nil {
		// Create default preferences
		return CreateUserPreference(user)
	}
	return &preference, nil
}
