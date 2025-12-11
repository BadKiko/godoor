package serializers

import (
	"godoor/models"
	"time"
)

// UnitSerializer matching Tandoor UnitSerializer
type UnitSerializer struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	PluralName  string `json:"plural_name"`
	Description string `json:"description"`
}

// UnitListResponse represents paginated unit response
type UnitListResponse struct {
	Count     int              `json:"count"`
	Next      *string          `json:"next"`
	Previous  *string          `json:"previous"`
	Results   []UnitSerializer `json:"results"`
	Timestamp string           `json:"timestamp"`
}

// SerializeUnit converts Unit model to serializer
func SerializeUnit(unit *models.Unit) *UnitSerializer {
	if unit == nil {
		return nil
	}
	pluralName := ""
	if unit.PluralName != nil {
		pluralName = *unit.PluralName
	}
	return &UnitSerializer{
		ID:          unit.ID,
		Name:        unit.Name,
		PluralName:  pluralName,
		Description: unit.Description,
	}
}

// SerializeUnits converts slice of Unit models to paginated response
func SerializeUnits(units []models.Unit, totalCount int) UnitListResponse {
	result := make([]UnitSerializer, len(units))
	for i, unit := range units {
		serialized := SerializeUnit(&unit)
		if serialized != nil {
			result[i] = *serialized
		}
	}
	return UnitListResponse{
		Count:     totalCount,
		Next:      nil, // TODO: implement pagination
		Previous:  nil, // TODO: implement pagination
		Results:   result,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
