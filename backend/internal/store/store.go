package store

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides typed access to PostgreSQL persistence.
type Store struct {
	db *sql.DB
}

// New creates a Store backed by the given database connection.
func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// Ping verifies the database connection is alive.
func (s *Store) Ping(ctx context.Context) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store: database not configured")
	}
	return s.db.PingContext(ctx)
}
