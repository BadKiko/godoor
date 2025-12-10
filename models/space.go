package models

// stringPtr creates a string pointer
func stringPtr(s string) *string {
	return &s
}

// Space model matching Tandoor structure
type Space struct {
	BaseModel
	Name        string    `json:"name" gorm:"default:'Default'"`
	CreatedByID uint      `json:"-" gorm:"column:created_by_id"`
	CreatedBy   User      `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID;references:ID"`
	Message     string    `json:"message" gorm:"default:''"`
}

// UserSpace represents the relationship between users and spaces
type UserSpace struct {
	BaseModel
	UserID      uint      `json:"-" gorm:"not null"`
	User        User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	SpaceID     uint      `json:"-" gorm:"not null"`
	Space       Space     `json:"space" gorm:"foreignKey:SpaceID;references:ID"`
	Active      bool      `json:"active" gorm:"default:false"`
	InternalNote *string  `json:"internal_note"`
}

// GetActiveSpace returns the active space for a user
func (u *User) GetActiveSpace() *Space {
	var userSpace UserSpace
	if err := DB.Where("user_id = ? AND active = ?", u.ID, true).Preload("Space").First(&userSpace).Error; err != nil {
		return nil
	}
	return &userSpace.Space
}

// CreateSpaceForUser creates a new space for user
func CreateSpaceForUser(user *User, name *string) (*UserSpace, error) {
	spaceName := "Default"
	if name != nil && *name != "" {
		spaceName = *name
	}

	space := Space{
		Name:        spaceName,
		CreatedByID: user.ID,
		CreatedBy:   *user,
	}

	if err := DB.Create(&space).Error; err != nil {
		return nil, err
	}

	userSpace := UserSpace{
		UserID:  user.ID,
		User:    *user,
		SpaceID: space.ID,
		Space:   space,
		Active:  true,
	}

	if err := DB.Create(&userSpace).Error; err != nil {
		return nil, err
	}

	// Create default meal types for new space
	defaultMealTypes := []struct {
		name  string
		order int
		time  *string
		color *string
	}{
		{"Breakfast", 1, stringPtr("08:00:00"), stringPtr("#FF6B35")},
		{"Lunch", 2, stringPtr("12:00:00"), stringPtr("#F7931E")},
		{"Dinner", 3, stringPtr("18:00:00"), stringPtr("#FFD23F")},
		{"Snack", 4, nil, stringPtr("#06FFA5")},
	}

	for _, mt := range defaultMealTypes {
		mealType := MealType{
			Name:      mt.name,
			Order:     mt.order,
			Time:      mt.time,
			Color:     mt.color,
			Default:   mt.order == 1, // Make breakfast default
			CreatedBy: *user,
			Space:     space,
		}

		if err := DB.Create(&mealType).Error; err != nil {
			// Log error but don't fail space creation
			continue
		}
	}

	return &userSpace, nil
}

// CreateSpace creates a new space
func CreateSpace(createdBy *User, name, message string) (*Space, error) {
	space := Space{
		Name:        name,
		CreatedByID: createdBy.ID,
		CreatedBy:   *createdBy,
		Message:     message,
	}

	if err := DB.Create(&space).Error; err != nil {
		return nil, err
	}

	return &space, nil
}
