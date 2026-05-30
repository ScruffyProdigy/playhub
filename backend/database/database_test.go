package database

import (
	"os"
	"testing"
)

func TestInitWithMigrationsKeepsConnectionOpen(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set")
	}

	if err := InitWithMigrations(); err != nil {
		t.Fatalf("InitWithMigrations() error: %v", err)
	}
	t.Cleanup(func() { _ = Close() })

	if err := GetDB().Ping(); err != nil {
		t.Fatalf("GetDB().Ping() after migrations: %v", err)
	}
}
