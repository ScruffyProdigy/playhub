package graph

import (
	"context"
	"sync"
	"testing"

	gqlclient "github.com/99designs/gqlgen/client"
	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/auth"
	"github.com/scruffyprodigy/playhub/internal/gameclient"
	"github.com/scruffyprodigy/playhub/internal/pubsub"
)

type syncProvisioner struct {
	mu    sync.Mutex
	calls []gameclient.Assignment
}

func (p *syncProvisioner) ProvisionMatch(_ context.Context, _, _ string, assignment gameclient.Assignment) error {
	p.mu.Lock()
	p.calls = append(p.calls, assignment)
	p.mu.Unlock()
	return nil
}

func (p *syncProvisioner) lastAssignment() gameclient.Assignment {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.calls) == 0 {
		return gameclient.Assignment{}
	}
	return p.calls[len(p.calls)-1]
}

func TestJoinGameQueueProvisionsMatchOnGameServer(t *testing.T) {
	env := newQueueIntegrationEnv(t)
	cleaner := env.newCleaner(t)
	ctx := context.Background()
	clearDemoGameWaitingQueue(t, env.Store)

	provisioner := &syncProvisioner{}
	env.resolverWithProvisioner(provisioner)

	gameID, _ := uuid.Parse(demoQuickMatchGameID)
	if err := env.Store.SetGameHandoffURLsForTest(ctx, gameID, "http://localhost:5174", "http://127.0.0.1:1"); err != nil {
		t.Fatalf("set urls: %v", err)
	}

	_, cookieA := createTestUserSession(t, ctx, env, cleaner)
	_, cookieB := createTestUserSession(t, ctx, env, cleaner)

	joinQuery := `mutation Join($id: ID!) { joinGame(gameId: $id) { queued queuedCount } }`
	vars := map[string]any{"id": demoQuickMatchGameID}

	postGraphQL(t, env.Handler, joinQuery, vars, cookieA)
	postGraphQL(t, env.Handler, joinQuery, vars, cookieB)

	if n := len(provisioner.calls); n != 1 {
		t.Fatalf("expected exactly 1 provision call, got %d", n)
	}
	a := provisioner.lastAssignment()
	if a.ExternalMatchID == "" || len(a.Seats) != 2 {
		t.Fatalf("expected provisioned duel with 2 seats, got %+v", a)
	}
	if a.GameMode != "duel" {
		t.Fatalf("gameMode = %q, want duel", a.GameMode)
	}
}

func (env *queueIntegrationEnv) resolverWithProvisioner(p MatchProvisioner) {
	authService, err := auth.NewService(env.Store, env.Signer)
	if err != nil {
		panic(err)
	}
	resolver := NewResolver(env.Store, authService, pubsub.NewMemory(), "http://localhost:5174")
	resolver.GameProvisioner = p
	gqlHandler := auth.Middleware(env.Signer, NewGraphQLServer(env.Signer, resolver))
	env.Handler = gqlHandler
	env.Client = gqlclient.New(gqlHandler)
}
