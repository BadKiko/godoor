package serializers

import (
	"godoor/models"
	"time"
)

// UserSpaceSerializer matching Tandoor UserSpaceSerializer
type UserSpaceSerializer struct {
	ID           uint           `json:"id"`
	User         UserSerializer `json:"user"`
	Space        SpaceSerializer `json:"space"`
	Groups       []interface{}  `json:"groups"` // TODO: implement groups
	Active       bool           `json:"active"`
	InternalNote *string        `json:"internal_note"`
	InviteLink   *interface{}   `json:"invite_link"` // TODO: implement invite links
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// SerializeUserSpace converts UserSpace model to UserSpaceSerializer
func SerializeUserSpace(userSpace *models.UserSpace) UserSpaceSerializer {
	return UserSpaceSerializer{
		ID:           userSpace.ID,
		User:         SerializeUser(&userSpace.User),
		Space:        SerializeSpace(&userSpace.Space),
		Groups:       []interface{}{}, // TODO: implement groups
		Active:       userSpace.Active,
		InternalNote: userSpace.InternalNote,
		InviteLink:   nil, // TODO: implement invite links
		CreatedAt:    userSpace.CreatedAt,
		UpdatedAt:    userSpace.UpdatedAt,
	}
}

// SerializeUserSpaces converts slice of UserSpace models to slice of UserSpaceSerializers
func SerializeUserSpaces(userSpaces []models.UserSpace) []UserSpaceSerializer {
	result := make([]UserSpaceSerializer, len(userSpaces))
	for i, userSpace := range userSpaces {
		result[i] = SerializeUserSpace(&userSpace)
	}
	return result
}
