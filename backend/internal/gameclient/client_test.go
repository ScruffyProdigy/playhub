package gameclient

import "testing"

func TestServiceAuthHeader(t *testing.T) {
	if got := serviceAuthHeader(""); got != "" {
		t.Fatalf("empty: %q", got)
	}
	if got := serviceAuthHeader("secret"); got != "Bearer secret" {
		t.Fatalf("raw: %q", got)
	}
	if got := serviceAuthHeader("Bearer already"); got != "Bearer already" {
		t.Fatalf("bearer: %q", got)
	}
}
