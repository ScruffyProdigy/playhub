package testdb

import (
	"errors"
	"net/url"
	"os"
	"strings"
	"testing"
)

const testDBName = "playhub_test"

// RequireURL returns DATABASE_URL for integration tests. It skips unless the database
// name is playhub_test, so go test does not mutate the dev playhub database by accident.
// Override with ALLOW_TESTS_ON_DEV_DB=1 when intentionally testing against playhub.
func RequireURL(t *testing.T) string {
	t.Helper()

	raw := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if raw == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}
	if os.Getenv("ALLOW_TESTS_ON_DEV_DB") == "1" {
		return raw
	}
	dbName, err := DatabaseName(raw)
	if err != nil {
		t.Fatalf("parse DATABASE_URL: %v", err)
	}
	if dbName != testDBName {
		t.Skip("integration tests require database " + testDBName +
			" (run ./scripts/test-backend.sh or: export DATABASE_URL=$(./scripts/db.sh test-url))")
	}
	return raw
}

// DatabaseName returns the database segment from a postgres connection URL.
func DatabaseName(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", err
	}
	name := strings.TrimPrefix(u.Path, "/")
	if idx := strings.Index(name, "?"); idx >= 0 {
		name = name[:idx]
	}
	if name == "" {
		return "", errors.New("database name is empty")
	}
	return name, nil
}

// IsTestDatabase reports whether raw points at the integration-test database.
func IsTestDatabase(raw string) bool {
	name, err := DatabaseName(raw)
	return err == nil && name == testDBName
}
