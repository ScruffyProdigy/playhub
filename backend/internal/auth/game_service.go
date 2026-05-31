package auth

import (
	"context"
	"crypto/subtle"
	"os"
	"strings"
)

type contextKeyGameService struct{}

// GameServiceTokenFromEnv is Lobby's server-to-server secret (provision body + player GraphQL).
func GameServiceTokenFromEnv() string {
	return strings.TrimSpace(os.Getenv("LOBBY_GAME_SERVICE_TOKEN"))
}

// MatchesGameServiceToken reports whether token equals the configured service token.
func MatchesGameServiceToken(token string) bool {
	expected := GameServiceTokenFromEnv()
	token = strings.TrimSpace(token)
	if expected == "" || token == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(token), []byte(expected)) == 1
}

// WithGameServiceAuth marks the context as authenticated via LOBBY_GAME_SERVICE_TOKEN.
func WithGameServiceAuth(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyGameService{}, true)
}

// GameServiceFromContext reports whether the request used the game service token.
func GameServiceFromContext(ctx context.Context) bool {
	ok, _ := ctx.Value(contextKeyGameService{}).(bool)
	return ok
}
