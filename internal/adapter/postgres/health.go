package postgres

import (
	"context"

	"gorm.io/gorm"
)

// Health is the readiness checker handed to the health use case.
type Health struct {
	db *gorm.DB
}

func NewHealth(db *gorm.DB) *Health { return &Health{db: db} }

func (h *Health) Ping(ctx context.Context) error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
