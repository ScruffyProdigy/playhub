package migrate

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Migrator handles database migrations
type Migrator struct {
	migrator *migrate.Migrate
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB) (*Migrator, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Get the migrations directory path
	migrationsPath := filepath.Join("migrations")
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		// Try relative to the current working directory
		migrationsPath = filepath.Join("backend", "migrations")
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &Migrator{migrator: m}, nil
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	log.Println("Running database migrations...")

	if err := m.migrator.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to run")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	log.Println("Rolling back last migration...")

	if err := m.migrator.Steps(-1); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Migration rollback completed successfully")
	return nil
}

// Steps runs a specific number of migrations (positive for up, negative for down)
func (m *Migrator) Steps(steps int) error {
	log.Printf("Running %d migration steps...", steps)

	if err := m.migrator.Steps(steps); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No migrations to run")
			return nil
		}
		return fmt.Errorf("failed to run migration steps: %w", err)
	}

	log.Printf("Migration steps completed successfully")
	return nil
}

// Version returns the current migration version
func (m *Migrator) Version() (uint, bool, error) {
	version, dirty, err := m.migrator.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}
	return version, dirty, nil
}

// Force sets the migration version (useful for fixing dirty migrations)
func (m *Migrator) Force(version int) error {
	log.Printf("Forcing migration version to %d...", version)

	if err := m.migrator.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	log.Printf("Migration version forced to %d", version)
	return nil
}

// Close closes the migrator
func (m *Migrator) Close() error {
	sourceErr, dbErr := m.migrator.Close()
	if sourceErr != nil {
		return fmt.Errorf("failed to close migration source: %w", sourceErr)
	}
	if dbErr != nil {
		return fmt.Errorf("failed to close migration database: %w", dbErr)
	}
	return nil
}
