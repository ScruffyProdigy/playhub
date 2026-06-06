package seattemplate

import "testing"

func TestValidateDistinctDisplayNames(t *testing.T) {
	t.Parallel()
	if err := ValidateDistinctDisplayNames([]PathSpec{
		{QueuePath: "DPS", DisplayName: "Damage"},
		{QueuePath: "Tank", DisplayName: "Tank"},
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := ValidateDistinctDisplayNames([]PathSpec{
		{QueuePath: "DPS", DisplayName: "Role"},
		{QueuePath: "Tank", DisplayName: "role"},
	})
	if err == nil {
		t.Fatal("expected duplicate displayName error")
	}
}
