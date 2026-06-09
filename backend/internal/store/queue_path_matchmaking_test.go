package store

import (
	"testing"
)

func TestMatchSeatsFromTemplateWordHunt(t *testing.T) {
	seats := []GameModeSeat{
		{SeatKey: "ClueGiver-Red", QueuePath: strPtr("ClueGiver"), SortOrder: 0},
		{SeatKey: "ClueGiver-Blue", QueuePath: strPtr("ClueGiver"), SortOrder: 1},
		{SeatKey: "ClueGiver-Green", QueuePath: strPtr("ClueGiver"), SortOrder: 2},
		{SeatKey: "Guesser-1", QueuePath: strPtr("Guesser"), SortOrder: 3},
		{SeatKey: "Guesser-2", QueuePath: strPtr("Guesser"), SortOrder: 4},
		{SeatKey: "Guesser-3", QueuePath: strPtr("Guesser"), SortOrder: 5},
		{SeatKey: "Guesser-4", QueuePath: strPtr("Guesser"), SortOrder: 6},
		{SeatKey: "Guesser-5", QueuePath: strPtr("Guesser"), SortOrder: 7},
		{SeatKey: "Guesser-6", QueuePath: strPtr("Guesser"), SortOrder: 8},
	}
	template := []byte(`{
		"ClueGiver":{"displayName":"Clue Giver","name":["Red","Blue","Green"],"min":2,"max":3,"sizeForQueue":2},
		"Guesser":{"count":6,"min":2,"max":6,"sizeForQueue":4}
	}`)
	matchSeats, err := matchSeatsFromTemplate(seats, template)
	if err != nil {
		t.Fatal(err)
	}
	if len(matchSeats) != 6 {
		t.Fatalf("got %d match seats, want 6", len(matchSeats))
	}
	keys := make([]string, len(matchSeats))
	for i, seat := range matchSeats {
		keys[i] = seat.SeatKey
	}
	want := []string{"ClueGiver-Red", "ClueGiver-Blue", "Guesser-1", "Guesser-2", "Guesser-3", "Guesser-4"}
	for i, key := range want {
		if keys[i] != key {
			t.Fatalf("seat[%d] = %q, want %q (all=%v)", i, keys[i], key, keys)
		}
	}
}

func TestValidateJoinQueuePath(t *testing.T) {
	dps := "DPS"
	seats := []GameModeSeat{{SeatKey: "1"}}
	if err := validateJoinQueuePath(seats, "DPS"); err == nil {
		t.Fatal("expected error for path on single-bucket mode")
	}

	compSeats := []GameModeSeat{{SeatKey: "Team-1-DPS-1", QueuePath: &dps}}
	if err := validateJoinQueuePath(compSeats, ""); err == nil {
		t.Fatal("expected error for missing path on composition mode")
	}
	if err := validateJoinQueuePath(compSeats, "Tank"); err == nil {
		t.Fatal("expected error for invalid path")
	}
	if err := validateJoinQueuePath(compSeats, "DPS"); err != nil {
		t.Fatalf("valid path: %v", err)
	}
}

func TestQueuePathColumnValueUsesEmptyStringForFifo(t *testing.T) {
	if got := queuePathColumnValue(""); got != "" {
		t.Fatalf("queuePathColumnValue(\"\") = %q, want empty string", got)
	}
	if got := queuePathColumnValue("  DPS  "); got != "DPS" {
		t.Fatalf("queuePathColumnValue trimmed = %q, want DPS", got)
	}
	if nullQueuePathColumn("") != nil {
		t.Fatal("nullQueuePathColumn(\"\") should stay nil for nullable game_queues.queue_path")
	}
}

func strPtr(s string) *string { return &s }
