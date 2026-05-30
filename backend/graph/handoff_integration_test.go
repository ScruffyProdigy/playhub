package graph

import (
	"context"
	"sync"
	"testing"

	"github.com/scruffyprodigy/playhub/internal/gameclient"
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
	env.resolverWithProvisioner(t, provisioner)

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

func (env *queueIntegrationEnv) resolverWithProvisioner(t *testing.T, p MatchProvisioner) {
	t.Helper()
	env.resolver.GameProvisioner = p
	env.rebuildHTTPServer(t)
}
