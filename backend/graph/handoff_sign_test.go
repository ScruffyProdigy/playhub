package graph

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/gameclient"
	"github.com/scruffyprodigy/playhub/internal/store"
)

type failingProvisioner struct{}

func (failingProvisioner) ProvisionMatch(context.Context, gameclient.ProvisionRequest) error {
	return errors.New("game server unavailable")
}

func TestSignLaunchURLMintsWhenProvisionFails(t *testing.T) {
	env := newQueueIntegrationEnv(t)
	cleaner := env.newCleaner(t)
	ctx := context.Background()
	clearDemoQueue(t, env.Store)

	t.Setenv("LOBBY_ISSUER_URL", "http://localhost:8080")
	t.Setenv("LOBBY_PUBLIC_URL", "http://localhost:5173")
	t.Setenv("LOBBY_GAME_TOKEN_PEPPER", "test-pepper")

	resolver := env.resolver
	env.resolverWithProvisioner(t, failingProvisioner{})

	userA, err := env.Store.CreateUser(ctx, store.CreateUserParams{
		Email: "launch-a-" + uuid.NewString() + "@example.com",
	})
	if err != nil {
		t.Fatalf("CreateUser A: %v", err)
	}
	cleaner.TrackUser(userA.ID)
	userB, err := env.Store.CreateUser(ctx, store.CreateUserParams{
		Email: "launch-b-" + uuid.NewString() + "@example.com",
	})
	if err != nil {
		t.Fatalf("CreateUser B: %v", err)
	}
	cleaner.TrackUser(userB.ID)

	queueID := uuid.MustParse(store.DemoDefaultQueueIDStr)
	if _, err := env.Store.JoinModeQueue(ctx, queueID, userA.ID, "", nil); err != nil {
		t.Fatalf("join A: %v", err)
	}
	result, err := env.Store.JoinModeQueue(ctx, queueID, userB.ID, "", nil)
	if err != nil {
		t.Fatalf("join B: %v", err)
	}
	if result.SessionID == nil {
		t.Fatal("expected matched session")
	}

	game, err := env.Store.GetGameByID(ctx, result.GameID)
	if err != nil {
		t.Fatalf("GetGameByID: %v", err)
	}
	playURL := "https://play.example.com"
	game.PlayURL = &playURL
	apiURL := "https://api.example.com"
	game.APIBaseURL = &apiURL

	launch, err := resolver.signLaunchURL(ctx, game, *result.SessionID, userB.ID)
	if err != nil {
		t.Fatalf("signLaunchURL: %v", err)
	}
	if launch == "" {
		t.Fatal("expected launch url despite provision failure")
	}
	if !strings.Contains(launch, "match=") || !strings.Contains(launch, "token=") {
		t.Fatalf("launch url = %q", launch)
	}
}
