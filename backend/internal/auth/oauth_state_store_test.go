package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestOAuthStateStoreRoundTrip(t *testing.T) {
	store := newMemoryOAuthStateStore()
	userID := uuid.New()
	state := OAuthState{
		Provider:     "discord",
		Mode:         OAuthModeLink,
		UserID:       userID,
		ConfirmMerge: true,
	}

	id, err := store.Save(context.Background(), state, "", time.Minute)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := uuid.Parse(id); err != nil {
		t.Fatalf("state id is not uuid: %q", id)
	}

	loaded, verifier, err := store.Load(context.Background(), id)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Provider != "discord" || loaded.Mode != OAuthModeLink || loaded.UserID != userID || !loaded.ConfirmMerge {
		t.Fatalf("loaded state = %#v", loaded)
	}
	if verifier != "" {
		t.Fatalf("verifier = %q", verifier)
	}

	if err := store.Delete(context.Background(), id); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, _, err := store.Load(context.Background(), id); err == nil {
		t.Fatal("expected load after delete to fail")
	}
}

func TestPKCEChallenge(t *testing.T) {
	verifier, err := generatePKCEVerifier()
	if err != nil {
		t.Fatalf("generatePKCEVerifier: %v", err)
	}
	if len(verifier) < 43 {
		t.Fatalf("verifier too short: %q", verifier)
	}
	challenge := pkceCodeChallenge(verifier)
	if challenge == "" || challenge == verifier {
		t.Fatalf("unexpected challenge: %q", challenge)
	}
}
