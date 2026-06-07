package tarot

import (
	"testing"
)

func TestDrawReturnsFiveDistinctCards(t *testing.T) {
	draw, err := Draw()
	if err != nil {
		t.Fatalf("Draw: %v", err)
	}
	if len(draw) != drawCount {
		t.Fatalf("expected %d cards, got %d", drawCount, len(draw))
	}
	seen := make(map[int]bool)
	for _, idx := range draw {
		if idx < 0 || idx >= cardCount {
			t.Fatalf("card index out of range: %d", idx)
		}
		if seen[idx] {
			t.Fatal("duplicate card drawn")
		}
		seen[idx] = true
	}
}

func TestCardName(t *testing.T) {
	if got, _ := CardName(0); got != "The Fool" {
		t.Fatalf("got %q", got)
	}
	if got, _ := CardName(21); got != "The World" {
		t.Fatalf("got %q", got)
	}
	if _, err := CardName(22); err == nil {
		t.Fatal("expected error for invalid index")
	}
}

func TestJourneySlotsCount(t *testing.T) {
	if len(JourneySlots) != drawCount {
		t.Fatalf("expected %d slots, got %d", drawCount, len(JourneySlots))
	}
}
