package serializers

import (
	"godoor/models"
)

// UserSerializer matching Tandoor UserSerializer
type UserSerializer struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
	IsStaff     bool   `json:"is_staff"`
	IsSuperuser bool   `json:"is_superuser"`
	IsActive    bool   `json:"is_active"`
}

// SerializeUser converts User model to UserSerializer
func SerializeUser(user *models.User) UserSerializer {
	return UserSerializer{
		ID:          user.ID,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DisplayName: user.GetUserDisplayName(),
		IsStaff:     user.IsStaff,
		IsSuperuser: user.IsSuperuser,
		IsActive:    user.IsActive,
	}
}

// SerializeUsers converts slice of User models to slice of UserSerializers
func SerializeUsers(users []models.User) []UserSerializer {
	result := make([]UserSerializer, len(users))
	for i, user := range users {
		result[i] = SerializeUser(&user)
	}
	return result
}
