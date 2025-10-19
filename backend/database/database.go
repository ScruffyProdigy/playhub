package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/scruffyprodigy/playhub/internal/migrate"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// Init initializes the database connection
func Init() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return nil
}

// InitWithMigrations initializes the database connection and runs migrations
func InitWithMigrations() error {
	// Initialize database connection
	if err := Init(); err != nil {
		return err
	}

	// Run migrations
	migrator, err := migrate.NewMigrator(DB)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	// Run all pending migrations
	if err := migrator.Up(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}
