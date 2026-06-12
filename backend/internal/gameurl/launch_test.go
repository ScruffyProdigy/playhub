package gameurl

import "testing"

func TestAttachSeatToken(t *testing.T) {
	got, err := AttachSeatToken("https://play.example.com/?match=abc&seat=1", "jwt.here")
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://play.example.com/?match=abc&seat=1&token=jwt.here" {
		t.Fatalf("got %q", got)
	}
}

func TestAttachSeatTokenPathStyle(t *testing.T) {
	got, err := AttachSeatToken("https://play.example.com/m/abc", "jwt.here")
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://play.example.com/m/abc?token=jwt.here" {
		t.Fatalf("got %q", got)
	}
}
