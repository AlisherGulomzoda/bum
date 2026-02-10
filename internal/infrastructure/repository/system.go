package repository

import (
	"context"
	"fmt"

	"bum-service/pkg/postgres"
	"bum-service/pkg/transaction"
)

// System is system repository.
type System struct {
	db      postgres.DB
	session func(context.Context) postgres.DB
}

// NewSystem creates a new system repository instance.
func NewSystem(db postgres.DB, session transaction.SessionDB) *System {
	return &System{
		db:      db,
		session: session.DB,
	}
}

// Ping checks whether the system is healthy.
func (s *System) Ping(_ context.Context) error {
	err := s.db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
