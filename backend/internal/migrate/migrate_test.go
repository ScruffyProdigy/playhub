package migrate

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestMigrator(t *testing.T) {
	// Skip if no database URL is provided
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping migration test")
	}

	// Connect to test database
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Create migrator
	migrator, err := NewMigrator(db)
	if err != nil {
		t.Fatalf("Failed to create migrator: %v", err)
	}
	defer migrator.Close()

	// Test version check
	version, dirty, err := migrator.Version()
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}
	t.Logf("Current version: %d, dirty: %v", version, dirty)

	// Test running migrations up (should be idempotent)
	if err := migrator.Up(); err != nil {
		t.Fatalf("Failed to run migrations up: %v", err)
	}

	// Check version again
	version, dirty, err = migrator.Version()
	if err != nil {
		t.Fatalf("Failed to get version after up: %v", err)
	}
	t.Logf("Version after up: %d, dirty: %v", version, dirty)

	// Verify that the users table exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to check if users table exists: %v", err)
	}
	if !exists {
		t.Error("Users table should exist after migration")
	}

	// Verify that the magic_links table exists
	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'magic_links')").Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to check if magic_links table exists: %v", err)
	}
	if !exists {
		t.Error("Magic_links table should exist after migration")
	}
}
