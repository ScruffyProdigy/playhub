package auth

import "testing"

func TestMatchesGameServiceToken(t *testing.T) {
	t.Setenv("LOBBY_GAME_SERVICE_TOKEN", "secret-token")

	if MatchesGameServiceToken("") {
		t.Fatal("empty token should not match")
	}
	if MatchesGameServiceToken("wrong") {
		t.Fatal("wrong token should not match")
	}
	if !MatchesGameServiceToken("secret-token") {
		t.Fatal("expected match")
	}
}
