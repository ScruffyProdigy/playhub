package auth

import (
	"context"
	"os"
	"strings"
)

type contextKeyGameService struct{}
type contextKeyGameServiceGameID struct{}

// GameServiceTokenFromEnv is Lobby's legacy global server-to-server secret (dev fallback).
func GameServiceTokenFromEnv() string {
	return strings.TrimSpace(os.Getenv("LOBBY_GAME_SERVICE_TOKEN"))
}

// GameServiceAuthEnabled reports whether game servers must authenticate for player GraphQL.
func GameServiceAuthEnabled() bool {
	if strings.TrimSpace(os.Getenv("LOBBY_GAME_TOKEN_PEPPER")) != "" {
		return true
	}
	return GameServiceTokenFromEnv() != ""
}

// WithGameServiceAuth marks the context as authenticated via LOBBY_GAME_SERVICE_TOKEN.
func WithGameServiceAuth(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyGameService{}, true)
}

// WithGameServiceGameID scopes a per-game service token to one catalog game.
func WithGameServiceGameID(ctx context.Context, gameID string) context.Context {
	return context.WithValue(ctx, contextKeyGameServiceGameID{}, strings.TrimSpace(gameID))
}

// GameServiceFromContext reports whether the request used the game service token.
func GameServiceFromContext(ctx context.Context) bool {
	ok, _ := ctx.Value(contextKeyGameService{}).(bool)
	return ok
}

// GameServiceGameIDFromContext returns the game id from a v1.* service token, if present.
func GameServiceGameIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(contextKeyGameServiceGameID{}).(string)
	return id, ok && id != ""
}
