package store

import (
	"testing"

	"github.com/google/uuid"
)

func TestTryFormMatchSingleBucket(t *testing.T) {
	seats := []GameModeSeat{
		{SeatKey: "1"},
		{SeatKey: "2"},
	}
	waiting := []QueueEntry{
		{ID: uuid.New(), UserID: uuid.New()},
		{ID: uuid.New(), UserID: uuid.New()},
	}
	got, ok := tryFormMatch(seats, waiting)
	if !ok || len(got) != 2 {
		t.Fatalf("tryFormMatch = %v, ok=%v", got, ok)
	}
	if got[0].SeatKey != "1" || got[1].SeatKey != "2" {
		t.Fatalf("seat order = %+v", got)
	}
}

func TestTryFormMatchSkipsDuplicateUser(t *testing.T) {
	userID := uuid.New()
	otherID := uuid.New()
	seats := []GameModeSeat{
		{SeatKey: "1"},
		{SeatKey: "2"},
	}
	waiting := []QueueEntry{
		{ID: uuid.New(), UserID: userID},
		{ID: uuid.New(), UserID: userID},
		{ID: uuid.New(), UserID: otherID},
	}
	got, ok := tryFormMatch(seats, waiting)
	if !ok || len(got) != 2 {
		t.Fatalf("expected match, got %+v ok=%v", got, ok)
	}
	if got[0].Entry.UserID == got[1].Entry.UserID {
		t.Fatal("expected distinct users in match roster")
	}
	seen := map[uuid.UUID]struct{}{got[0].Entry.UserID: {}, got[1].Entry.UserID: {}}
	if _, ok := seen[userID]; !ok {
		t.Fatalf("expected user in match, got %+v", got)
	}
	if _, ok := seen[otherID]; !ok {
		t.Fatalf("expected other user in match, got %+v", got)
	}
}

func TestTryFormMatchCompositionRequiresMatchingPaths(t *testing.T) {
	dps := "DPS"
	tank := "Tank"
	seats := []GameModeSeat{
		{SeatKey: "Team-1-DPS-1", QueuePath: &dps},
		{SeatKey: "Team-1-Tank-1", QueuePath: &tank},
	}
	userA := uuid.New()
	userB := uuid.New()
	waiting := []QueueEntry{
		{ID: uuid.New(), UserID: userA, QueuePath: &dps},
		{ID: uuid.New(), UserID: userB, QueuePath: &dps},
	}
	if _, ok := tryFormMatch(seats, waiting); ok {
		t.Fatal("expected no match when tank seat has no tank waiter")
	}

	waiting = []QueueEntry{
		{ID: uuid.New(), UserID: userA, QueuePath: &dps},
		{ID: uuid.New(), UserID: userB, QueuePath: &tank},
	}
	got, ok := tryFormMatch(seats, waiting)
	if !ok || len(got) != 2 {
		t.Fatalf("expected match, got %+v ok=%v", got, ok)
	}
	if got[0].SeatKey != "Team-1-DPS-1" || got[1].SeatKey != "Team-1-Tank-1" {
		t.Fatalf("assignments = %+v", got)
	}
}

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

func strPtr(s string) *string { return &s }

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
