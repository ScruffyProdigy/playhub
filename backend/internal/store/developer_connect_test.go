package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestConnectOwnedGameAndRotateWebhook(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	owner, err := st.CreateUser(ctx, CreateUserParams{
		Email:       "owner-" + uuid.NewString() + "@example.com",
		DisplayName: "Owner",
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(owner.ID)

	slug := "dev-" + uuid.NewString()[:8]
	manifest := sampleManifest()
	manifest.SHA256Hash = uuid.NewString()

	registered, err := st.RegisterMyGame(ctx, RegisterMyGameParams{
		OwnerUserID:      owner.ID,
		Slug:             slug,
		Name:             "Dev Game",
		ShortDescription: "A short blurb",
		APIBaseURL:       "https://api.example.com/" + slug,
		ContactEmail:     "dev@example.com",
	}, nil, fmt.Errorf("unreachable"))
	if err != nil {
		t.Fatalf("RegisterMyGame draft: %v", err)
	}
	cleaner.TrackGame(registered.Game.ID)
	if registered.Connected || registered.Game.Visibility != GameVisibilityDraft {
		t.Fatalf("expected disconnected draft, got connected=%v visibility=%s", registered.Connected, registered.Game.Visibility)
	}

	connected, err := st.ConnectOwnedGame(ctx, registered.Game.ID, owner.ID, nil, manifest, nil)
	if err != nil {
		t.Fatalf("ConnectOwnedGame: %v", err)
	}
	if !connected.Connected || !connected.Changed {
		t.Fatalf("expected connected+changed, got %+v", connected)
	}
	if connected.Game.Visibility != GameVisibilityPrivateTesting {
		t.Fatalf("expected private_testing, got %s", connected.Game.Visibility)
	}

	newURL := "https://api.example.com/" + slug + "-v2"
	updatedManifest := sampleManifest()
	updatedManifest.SHA256Hash = uuid.NewString()
	updatedManifest.Modes[0].DisplayName = "Classic v2"
	reconnect, err := st.ConnectOwnedGame(ctx, registered.Game.ID, owner.ID, &newURL, updatedManifest, nil)
	if err != nil {
		t.Fatalf("ConnectOwnedGame new URL: %v", err)
	}
	if !reconnect.Connected {
		t.Fatal("expected reconnect success")
	}
	if reconnect.Game.APIBaseURL == nil || *reconnect.Game.APIBaseURL != newURL {
		t.Fatalf("expected apiBaseUrl %q, got %+v", newURL, reconnect.Game.APIBaseURL)
	}

	if _, err := st.db.ExecContext(ctx, `UPDATE games SET visibility = $2 WHERE id = $1`, registered.Game.ID, GameVisibilityPublic); err != nil {
		t.Fatalf("force public: %v", err)
	}
	badURL := "https://evil.example.com"
	_, err = st.ConnectOwnedGame(ctx, registered.Game.ID, owner.ID, &badURL, updatedManifest, nil)
	if err == nil {
		t.Fatal("expected apiBaseUrl change blocked for public game")
	}

	secret, game, err := st.RotateMyGameWebhookSecret(ctx, registered.Game.ID, owner.ID)
	if err != nil {
		t.Fatalf("RotateMyGameWebhookSecret: %v", err)
	}
	if secret == "" || game.WebhookSecret == nil || *game.WebhookSecret != secret {
		t.Fatalf("expected rotated secret, got secret=%q game=%+v", secret, game.WebhookSecret)
	}

	name := "Renamed Game"
	email := "new@example.com"
	site := "https://example.com"
	updated, err := st.UpdateMyGameMetadata(ctx, registered.Game.ID, owner.ID, UpdateMyGameMetadataParams{
		Name:         &name,
		ContactEmail: &email,
		WebsiteURL:   &site,
	})
	if err != nil {
		t.Fatalf("UpdateMyGameMetadata: %v", err)
	}
	if updated.Name != name {
		t.Fatalf("expected name %q, got %q", name, updated.Name)
	}
	if updated.ContactEmail == nil || *updated.ContactEmail != email {
		t.Fatalf("expected contact %q, got %+v", email, updated.ContactEmail)
	}
	if updated.WebsiteURL == nil || *updated.WebsiteURL != site {
		t.Fatalf("expected website %q, got %+v", site, updated.WebsiteURL)
	}
}
