package graph

import (
	"context"
	"sync"
	"testing"

	"github.com/scruffyprodigy/playhub/internal/gameclient"
)

type provisionCall struct {
	LobbyID    string
	Lobby      gameclient.LobbyInfo
	Assignment gameclient.Assignment
}

type syncProvisioner struct {
	mu    sync.Mutex
	calls []provisionCall
}

func (p *syncProvisioner) ProvisionMatch(_ context.Context, req gameclient.ProvisionRequest) error {
	p.mu.Lock()
	p.calls = append(p.calls, provisionCall{LobbyID: req.LobbyID, Lobby: req.Lobby, Assignment: req.Assignment})
	p.mu.Unlock()
	return nil
}

func (p *syncProvisioner) lastCall() provisionCall {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.calls) == 0 {
		return provisionCall{}
	}
	return p.calls[len(p.calls)-1]
}

func TestJoinQueueProvisionsMatchOnGameServer(t *testing.T) {
	env := newQueueIntegrationEnv(t)
	cleaner := env.newCleaner(t)
	ctx := context.Background()
	clearDemoQueue(t, env.Store)

	t.Setenv("LOBBY_ISSUER_URL", "http://localhost:8080")
	t.Setenv("LOBBY_PUBLIC_URL", "http://localhost:5173")

	provisioner := &syncProvisioner{}
	env.resolverWithProvisioner(t, provisioner)

	_, cookieA := createTestUserSession(t, ctx, env, cleaner)
	_, cookieB := createTestUserSession(t, ctx, env, cleaner)

	joinQuery := `mutation Join($id: ID!) { joinQueue(queueId: $id) { queued queuedCount } }`
	vars := map[string]any{"id": demoDefaultQueueID}

	postGraphQL(t, env.Handler, joinQuery, vars, cookieA)
	postGraphQL(t, env.Handler, joinQuery, vars, cookieB)

	if n := len(provisioner.calls); n != 1 {
		t.Fatalf("expected exactly 1 provision call, got %d", n)
	}
	call := provisioner.lastCall()
	a := call.Assignment
	if call.LobbyID != "http://localhost:8080" {
		t.Fatalf("lobbyId = %q, want http://localhost:8080", call.LobbyID)
	}
	if call.Lobby.ReturnURL != "http://localhost:5173" {
		t.Fatalf("lobby.returnUrl = %q", call.Lobby.ReturnURL)
	}
	if call.Lobby.GraphqlURL != "http://localhost:8080/graphql" {
		t.Fatalf("lobby.graphqlUrl = %q", call.Lobby.GraphqlURL)
	}
	if call.Lobby.ServiceToken != "" {
		t.Fatalf("lobby.serviceToken = %q, want empty when LOBBY_GAME_SERVICE_TOKEN unset", call.Lobby.ServiceToken)
	}
	if a.ExternalMatchID == "" || len(a.Seats) != 2 {
		t.Fatalf("expected provisioned duel with 2 seats, got %+v", a)
	}
	if a.GameMode != "duel" {
		t.Fatalf("gameMode = %q, want duel", a.GameMode)
	}
	if a.Seats[0].SeatKey != "a" || a.Seats[1].SeatKey != "b" {
		t.Fatalf("expected manifest seat keys a/b, got %+v", a.Seats)
	}
}

func TestJoinQueueProvisionsServiceTokenWhenConfigured(t *testing.T) {
	env := newQueueIntegrationEnv(t)
	cleaner := env.newCleaner(t)
	ctx := context.Background()
	clearDemoQueue(t, env.Store)

	t.Setenv("LOBBY_ISSUER_URL", "http://localhost:8080")
	t.Setenv("LOBBY_PUBLIC_URL", "http://localhost:5173")
	t.Setenv("LOBBY_GAME_SERVICE_TOKEN", "dev-lobby-svc-token")

	provisioner := &syncProvisioner{}
	env.resolverWithProvisioner(t, provisioner)

	_, cookieA := createTestUserSession(t, ctx, env, cleaner)
	_, cookieB := createTestUserSession(t, ctx, env, cleaner)

	joinQuery := `mutation Join($id: ID!) { joinQueue(queueId: $id) { queued queuedCount } }`
	vars := map[string]any{"id": demoDefaultQueueID}

	postGraphQL(t, env.Handler, joinQuery, vars, cookieA)
	postGraphQL(t, env.Handler, joinQuery, vars, cookieB)

	call := provisioner.lastCall()
	if call.Lobby.ServiceToken != "dev-lobby-svc-token" {
		t.Fatalf("lobby.serviceToken = %q, want dev-lobby-svc-token", call.Lobby.ServiceToken)
	}
}
