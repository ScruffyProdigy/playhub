package seattemplate

import (
	"encoding/json"
	"testing"
)

func mustExpand(t *testing.T, template string) []Leaf {
	t.Helper()
	leaves, err := Expand(json.RawMessage(template))
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	return leaves
}

func seatKeys(leaves []Leaf) []string {
	out := make([]string, len(leaves))
	for i, leaf := range leaves {
		out[i] = leaf.SeatKey
	}
	return out
}

func TestExpandDuelCountTwo(t *testing.T) {
	leaves := mustExpand(t, `{"count":2}`)
	if got := seatKeys(leaves); got[0] != "1" || got[1] != "2" {
		t.Fatalf("seat keys = %v, want [1 2]", got)
	}
	if leaves[0].QueuePath != "" || leaves[1].QueuePath != "" {
		t.Fatalf("expected default queue path, got %+v", leaves)
	}
	if leaves[0].AffinityKey != "" {
		t.Fatalf("affinity should be empty in Phase A, got %q", leaves[0].AffinityKey)
	}
}

func TestExpandThreeByThree(t *testing.T) {
	leaves := mustExpand(t, `{"Team":{"count":2,"Seat":{"count":3}}}`)
	want := []string{
		"Team-1-Seat-1", "Team-1-Seat-2", "Team-1-Seat-3",
		"Team-2-Seat-1", "Team-2-Seat-2", "Team-2-Seat-3",
	}
	got := seatKeys(leaves)
	if len(got) != len(want) {
		t.Fatalf("got %d seats, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("seat[%d] = %q, want %q (all=%v)", i, got[i], want[i], got)
		}
	}
}

func TestExpandEmptyObjectVsCountOne(t *testing.T) {
	leaves := mustExpand(t, `{"Team":{"count":1,"Support":{}}}`)
	if leaves[0].SeatKey != "Team-1-Support" {
		t.Fatalf("{} seat = %q, want Team-1-Support", leaves[0].SeatKey)
	}

	leaves = mustExpand(t, `{"Team":{"count":1,"Support":{"count":1}}}`)
	if leaves[0].SeatKey != "Team-1-Support-1" {
		t.Fatalf("{count:1} seat = %q, want Team-1-Support-1", leaves[0].SeatKey)
	}
}

func TestExpandCompositionQueuePaths(t *testing.T) {
	leaves := mustExpand(t, `{
		"Team":{"count":2,
			"DPS":{"count":2},
			"Tank":{"count":2},
			"Support":{}
		}
	}`)
	if len(leaves) != 10 {
		t.Fatalf("expected 10 seats, got %d", len(leaves))
	}
	paths := map[string]int{}
	for _, leaf := range leaves {
		paths[leaf.QueuePath]++
	}
	if paths["DPS"] != 4 || paths["Tank"] != 4 || paths["Support"] != 2 {
		t.Fatalf("unexpected queue paths: %+v", paths)
	}
}

func TestExpandRejectsFlatSeats(t *testing.T) {
	_, err := Expand(json.RawMessage(`{"seats":[{"key":"a"}]}`))
	if err == nil {
		t.Fatal("expected error for flat seats[]")
	}
}

func TestExpandNameArray(t *testing.T) {
	leaves := mustExpand(t, `{"ClueGiver":{"name":["Red","Blue","Green"]}}`)
	want := []string{"ClueGiver-Red", "ClueGiver-Blue", "ClueGiver-Green"}
	got := seatKeys(leaves)
	if len(got) != len(want) {
		t.Fatalf("got %v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("seat[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestExpandWordHuntPartyTemplate(t *testing.T) {
	leaves := mustExpand(t, `{
		"ClueGiver":{"displayName":"Clue Giver","name":["Red","Blue","Green"],"min":2,"max":3,"sizeForQueue":2},
		"Guesser":{"count":6,"min":2,"max":6,"sizeForQueue":4}
	}`)
	want := []string{
		"ClueGiver-Red", "ClueGiver-Blue", "ClueGiver-Green",
		"Guesser-1", "Guesser-2", "Guesser-3", "Guesser-4", "Guesser-5", "Guesser-6",
	}
	got := seatKeys(leaves)
	if len(got) != len(want) {
		t.Fatalf("got %d seats %v, want %d", len(got), got, len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("seat[%d] = %q, want %q", i, got[i], want[i])
		}
	}
	paths := map[string]int{}
	for _, leaf := range leaves {
		paths[leaf.QueuePath]++
	}
	if paths["ClueGiver"] != 3 || paths["Guesser"] != 6 {
		t.Fatalf("unexpected queue paths: %+v", paths)
	}
}

func TestExpandSiegeOffenseDefense(t *testing.T) {
	leaves := mustExpand(t, `{
		"Offense":{"Wizard":{},"Warrior":{"count":2}},
		"Defense":{"Warrior":{"count":2}}
	}`)
	want := []string{
		"Defense-Warrior-1", "Defense-Warrior-2",
		"Offense-Warrior-1", "Offense-Warrior-2", "Offense-Wizard",
	}
	got := seatKeys(leaves)
	if len(got) != len(want) {
		t.Fatalf("got %v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("seat[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
