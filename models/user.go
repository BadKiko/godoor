package models

import (
	"godoor/utils"
	"gorm.io/gorm"
)

// User model matching Tandoor structure
type User struct {
	BaseModel
	Username   string `json:"username" gorm:"unique;not null"`
	Password   string `json:"-" gorm:"not null"` // Never serialize password
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	IsStaff    bool   `json:"is_staff" gorm:"default:false"`
	IsSuperuser bool  `json:"is_superuser" gorm:"default:false"`
	IsActive   bool   `json:"is_active" gorm:"default:true"`
}

// GetUserDisplayName returns display name for user (first_name + last_name or username)
func (u *User) GetUserDisplayName() string {
	if name := u.FirstName + " " + u.LastName; name != " " {
		return name
	}
	return u.Username
}

// BeforeCreate hook for GORM
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// CreateUser creates a new user with hashed password
func CreateUser(username, password, firstName, lastName string) (*User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := User{
		Username:  username,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
	}

	if err := DB.Create(&user).Error; err != nil {
		return nil, err
	}

	// Create default space for user
	_, err = CreateSpaceForUser(&user, nil)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CheckPassword verifies user password
func (u *User) CheckPassword(password string) bool {
	return utils.CheckPasswordHash(password, u.Password)
}
