package auth

import "testing"

func TestGameServiceAuthEnabled(t *testing.T) {
	t.Setenv("LOBBY_GAME_TOKEN_PEPPER", "")
	t.Setenv("LOBBY_GAME_SERVICE_TOKEN", "")
	if GameServiceAuthEnabled() {
		t.Fatal("expected auth disabled with no secrets")
	}

	t.Setenv("LOBBY_GAME_TOKEN_PEPPER", "pepper-only")
	if !GameServiceAuthEnabled() {
		t.Fatal("expected auth enabled with LOBBY_GAME_TOKEN_PEPPER")
	}

	t.Setenv("LOBBY_GAME_TOKEN_PEPPER", "")
	t.Setenv("LOBBY_GAME_SERVICE_TOKEN", "legacy-only")
	if !GameServiceAuthEnabled() {
		t.Fatal("expected auth enabled with LOBBY_GAME_SERVICE_TOKEN")
	}
}

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
