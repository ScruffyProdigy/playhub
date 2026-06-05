package store

import (
	"testing"

	"github.com/google/uuid"
)

func TestTryFormMatchFifo(t *testing.T) {
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

func TestValidateJoinQueuePath(t *testing.T) {
	dps := "DPS"
	seats := []GameModeSeat{{SeatKey: "1"}}
	if err := validateJoinQueuePath(seats, "DPS"); err == nil {
		t.Fatal("expected error for path on fifo mode")
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
