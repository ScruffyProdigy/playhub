package graph

import (
	"encoding/json"
	"testing"
)

func TestQueuePathDisplayName(t *testing.T) {
	template := json.RawMessage(`{
		"ClueGiver": {"displayName": "Clue Giver", "name": ["Red"], "min": 2, "max": 3, "sizeForQueue": 2},
		"Guesser": {"count": 6, "min": 2, "max": 6, "sizeForQueue": 4}
	}`)

	clue := "ClueGiver"
	got := queuePathDisplayName(template, &clue)
	if got == nil || *got != "Clue Giver" {
		t.Fatalf("ClueGiver label = %v, want Clue Giver", got)
	}

	guesser := "Guesser"
	got = queuePathDisplayName(template, &guesser)
	if got == nil || *got != "Guesser" {
		t.Fatalf("Guesser label = %v, want Guesser", got)
	}

	if queuePathDisplayName(template, nil) != nil {
		t.Fatal("expected nil for nil path")
	}
	empty := ""
	if queuePathDisplayName(template, &empty) != nil {
		t.Fatal("expected nil for empty path")
	}
}
