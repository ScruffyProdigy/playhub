package graph

import (
	"encoding/json"
	"testing"

	"github.com/scruffyprodigy/playhub/internal/store"
)

func TestSeatDisplayName(t *testing.T) {
	t.Parallel()
	template := json.RawMessage(`{
		"ClueGiver":{"displayName":"Clue Giver","name":["Red","Blue","Green"],"min":2,"max":3,"sizeForQueue":2},
		"Guesser":{"count":2,"min":2,"max":6,"sizeForQueue":4}
	}`)
	seats := []store.GameModeSeat{
		{SeatKey: "ClueGiver-Red", QueuePath: strPtr("ClueGiver")},
		{SeatKey: "Guesser-1", QueuePath: strPtr("Guesser")},
		{SeatKey: "1"},
	}

	if got := seatDisplayName(seats, "ClueGiver-Red", template); got != "Clue Giver · Red" {
		t.Fatalf("ClueGiver-Red = %q, want Clue Giver · Red", got)
	}
	if got := seatDisplayName(seats, "Guesser-1", template); got != "Guesser · 1" {
		t.Fatalf("Guesser-1 = %q, want Guesser · 1", got)
	}
	if got := seatDisplayName(seats, "1", template); got != "1" {
		t.Fatalf("fifo seat = %q, want 1", got)
	}
}

func strPtr(s string) *string {
	return &s
}
