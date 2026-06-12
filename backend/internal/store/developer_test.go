package store

import "testing"

func TestValidateGameSlug(t *testing.T) {
	if err := ValidateGameSlug("my-cool-game"); err != nil {
		t.Fatalf("expected valid slug: %v", err)
	}
	if err := ValidateGameSlug("Bad Slug"); err == nil {
		t.Fatal("expected invalid slug with spaces")
	}
	if err := ValidateGameSlug(""); err == nil {
		t.Fatal("expected empty slug to fail")
	}
}
