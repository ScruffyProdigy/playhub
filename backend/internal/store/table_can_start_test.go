package store

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestTableCanStartDuelWithTwoSeated(t *testing.T) {
	mode := &GameMode{
		MinPlayers:   2,
		MaxPlayers:   2,
		SeatTemplate: json.RawMessage(`{"count":2}`),
	}
	modeSeats := []GameModeSeat{
		{SeatKey: "1"},
		{SeatKey: "2"},
	}
	seated := []TableSeat{
		{SeatKey: "1", UserID: uuid.New()},
		{SeatKey: "2", UserID: uuid.New()},
	}

	ok, err := tableCanStart(mode, modeSeats, seated, mode.SeatTemplate)
	if err != nil {
		t.Fatalf("tableCanStart: %v", err)
	}
	if !ok {
		t.Fatal("expected duel table with 2 seated to be startable")
	}
}

func TestTableCanStartDuelWithOneSeated(t *testing.T) {
	mode := &GameMode{
		MinPlayers:   2,
		MaxPlayers:   2,
		SeatTemplate: json.RawMessage(`{"count":2}`),
	}
	modeSeats := []GameModeSeat{
		{SeatKey: "1"},
		{SeatKey: "2"},
	}
	seated := []TableSeat{
		{SeatKey: "1", UserID: uuid.New()},
	}

	ok, err := tableCanStart(mode, modeSeats, seated, mode.SeatTemplate)
	if err != nil {
		t.Fatalf("tableCanStart: %v", err)
	}
	if ok {
		t.Fatal("expected duel table with 1 seated to not be startable")
	}
}

func TestResolvePooledSeatKeyAutoAssignsNextOpen(t *testing.T) {
	modeSeats := []GameModeSeat{
		{SeatKey: "1", SortOrder: 0},
		{SeatKey: "2", SortOrder: 1},
	}
	seated := []TableSeat{{SeatKey: "1", UserID: uuid.New()}}

	got, err := resolvePooledSeatKey(modeSeats, seated, "1")
	if err != nil {
		t.Fatalf("resolvePooledSeatKey: %v", err)
	}
	if got != "2" {
		t.Fatalf("seat key = %q, want 2", got)
	}
}

func TestResolvePooledSeatKeyGroupFull(t *testing.T) {
	modeSeats := []GameModeSeat{
		{SeatKey: "1", SortOrder: 0},
		{SeatKey: "2", SortOrder: 1},
	}
	seated := []TableSeat{
		{SeatKey: "1", UserID: uuid.New()},
		{SeatKey: "2", UserID: uuid.New()},
	}

	_, err := resolvePooledSeatKey(modeSeats, seated, "1")
	if err == nil {
		t.Fatal("expected error when pooled group is full")
	}
}
