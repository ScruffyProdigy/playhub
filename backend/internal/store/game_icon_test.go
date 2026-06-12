package store

import "testing"

func TestDefaultGameIconURL(t *testing.T) {
	if got := DefaultGameIconURL(""); got != "/games/default.svg" {
		t.Fatalf("empty slug = %q", got)
	}
	if got := DefaultGameIconURL("word-hunt"); got != "/games/word-hunt-icon.png" {
		t.Fatalf("word-hunt = %q", got)
	}
	if got := DefaultGameIconURL("rock-paper-scissors-lizard-robot"); got != "/games/rpslr-icon.png" {
		t.Fatalf("rps slug override = %q", got)
	}
}

func TestDefaultGameHeroURL(t *testing.T) {
	if got := DefaultGameHeroURL(""); got != "/games/default-hero.svg" {
		t.Fatalf("empty slug = %q", got)
	}
	if got := DefaultGameHeroURL("word-hunt"); got != "/games/word-hunt-hero.jpg" {
		t.Fatalf("word-hunt = %q", got)
	}
	if got := DefaultGameHeroURL("rock-paper-scissors-lizard-robot"); got != "/games/rpslr-hero.jpg" {
		t.Fatalf("rps slug override = %q", got)
	}
}
