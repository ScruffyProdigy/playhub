package auth

import "testing"

func TestMagicLinkTokenFromURL(t *testing.T) {
	token, err := MagicLinkTokenFromURL("http://localhost:5173/auth/complete?token=abc-123")
	if err != nil {
		t.Fatalf("MagicLinkTokenFromURL: %v", err)
	}
	if token != "abc-123" {
		t.Fatalf("token = %q", token)
	}
}
