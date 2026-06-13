package integrationchecks

import "testing"

func TestLaunchURLsContainJWT(t *testing.T) {
	if launchURLsContainJWT(map[string]string{"u1": "https://play.example/?match=a&seat=1"}, "") {
		t.Fatal("expected clean URL")
	}
	if !launchURLsContainJWT(map[string]string{"u1": "https://play.example/?match=a&token=eyJhbG"}, "") {
		t.Fatal("expected token= in URL to fail")
	}
	if !launchURLsContainJWT(nil, "https://play.example/?match={matchId}&token={jwt}") {
		t.Fatal("expected template with token placeholder to fail")
	}
}
