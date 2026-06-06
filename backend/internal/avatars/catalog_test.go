package avatars

import "testing"

func TestStarterByKey(t *testing.T) {
	entry, ok := StarterByKey("Compass")
	if !ok || entry.Key != "compass" {
		t.Fatalf("expected compass, got %+v ok=%v", entry, ok)
	}
	_, ok = StarterByKey("unknown")
	if ok {
		t.Fatal("expected unknown key to miss")
	}
}

func TestPublicAssetURL(t *testing.T) {
	got := PublicAssetURL("https://joinquest.cc", "compass.png")
	want := "https://joinquest.cc/avatars/compass.png"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestResolveURLPrefersStored(t *testing.T) {
	stored := "https://cdn.example/avatar.png"
	got := ResolveURL("https://joinquest.cc", &stored, ptr("compass"))
	if got == nil || *got != stored {
		t.Fatalf("expected stored url, got %+v", got)
	}
}

func TestResolveURLFromKey(t *testing.T) {
	key := "beacon"
	got := ResolveURL("https://joinquest.cc", nil, &key)
	if got == nil || *got != "https://joinquest.cc/avatars/beacon.png" {
		t.Fatalf("unexpected url: %+v", got)
	}
}

func ptr(s string) *string { return &s }
