package testdb

import "testing"

func TestDatabaseName(t *testing.T) {
	tests := []struct {
		raw  string
		want string
	}{
		{"postgres://app:app-pass@127.0.0.1:5432/playhub_test?sslmode=disable", "playhub_test"},
		{"postgres://app:app-pass@127.0.0.1:5432/playhub?sslmode=disable", "playhub"},
	}
	for _, tc := range tests {
		got, err := DatabaseName(tc.raw)
		if err != nil {
			t.Fatalf("%q: %v", tc.raw, err)
		}
		if got != tc.want {
			t.Fatalf("%q: got %q want %q", tc.raw, got, tc.want)
		}
	}
}

func TestIsTestDatabase(t *testing.T) {
	if !IsTestDatabase("postgres://app:app-pass@127.0.0.1:5432/playhub_test?sslmode=disable") {
		t.Fatal("expected playhub_test")
	}
	if IsTestDatabase("postgres://app:app-pass@127.0.0.1:5432/playhub?sslmode=disable") {
		t.Fatal("expected not playhub_test")
	}
}
