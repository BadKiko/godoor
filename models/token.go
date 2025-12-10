package models

import (
	"time"
	"github.com/google/uuid"
)

// AccessToken model matching Tandoor OAuth2 token structure
type AccessToken struct {
	BaseModel
	UserID      uint      `json:"-" gorm:"not null"`
	User        User      `json:"-" gorm:"foreignKey:UserID;references:ID"`
	Token       string    `json:"token" gorm:"unique;not null"`
	Scope       string    `json:"scope" gorm:"default:'read write app'"`
	Expires     time.Time `json:"expires"`
}

// CreateTokenForUser creates a new access token for user (matching Tandoor logic)
func CreateTokenForUser(user *User) (*AccessToken, error) {
	// Check if user already has a valid token
	var existingToken AccessToken
	now := time.Now()
	expires := now.AddDate(5, 0, 0) // 5 years like in Tandoor

	err := DB.Where("user_id = ? AND expires > ?", user.ID, now).First(&existingToken).Error
	if err == nil && existingToken.Scope == "read write app" {
		return &existingToken, nil
	}

	// Generate new token
	token := "tda_" + uuid.New().String()

	accessToken := AccessToken{
		UserID:  user.ID,
		User:    *user,
		Token:   token,
		Scope:   "read write app",
		Expires: expires,
	}

	if err := DB.Create(&accessToken).Error; err != nil {
		return nil, err
	}

	return &accessToken, nil
}
