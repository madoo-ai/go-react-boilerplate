package entity

import "github.com/google/uuid"

// NewUUID returns a sortable UUIDv7. IDs are generated in the application so
// the database stays portable and tests can predict them when needed.
func NewUUID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.NewString()
	}
	return id.String()
}
