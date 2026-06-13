package developer

import "testing"

func TestValidateCatalogTags(t *testing.T) {
	if err := ValidateCatalogTags([]string{"competitive", "1v1"}); err != nil {
		t.Fatalf("expected valid tags: %v", err)
	}
	if err := ValidateCatalogTags([]string{"unknown"}); err == nil {
		t.Fatal("expected unknown tag to fail")
	}
}
