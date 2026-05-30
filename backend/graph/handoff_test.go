package graph

import (
	"testing"

	"github.com/scruffyprodigy/playhub/internal/pubsub"
	"github.com/scruffyprodigy/playhub/internal/store"
)

func TestResolvedAPIBaseURL(t *testing.T) {
	t.Setenv("GAME_API_BASE_URL", "http://game-api.test")

	r := NewResolver(nil, nil, pubsub.NewMemory(), "http://localhost:5174")

	game := &store.Game{Name: "Test"}
	if got := r.resolvedAPIBaseURL(game); got != "http://game-api.test" {
		t.Fatalf("resolvedAPIBaseURL = %q, want http://game-api.test", got)
	}

	api := "http://catalog-api.test"
	game.APIBaseURL = &api
	if got := r.resolvedAPIBaseURL(game); got != "http://catalog-api.test" {
		t.Fatalf("resolvedAPIBaseURL = %q, want catalog override", got)
	}
}

func TestBuildLaunchURL(t *testing.T) {
	got, err := buildLaunchURL("http://localhost:5174", "sess-1", "jwt.here")
	if err != nil {
		t.Fatal(err)
	}
	if got != "http://localhost:5174?match=sess-1&token=jwt.here" {
		t.Fatalf("got %q", got)
	}
}
