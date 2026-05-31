package auth

import (
	"net/url"
	"os"
	"strings"
)

const defaultLobbyIssuer = "http://localhost:8080"

// LobbyIssuer returns the canonical issuer URL for this Lobby deployment.
// Used as JWT iss on seat tokens and as lobbyId on game provision requests.
// Resolution order: LOBBY_ISSUER_URL, LOBBY_PUBLIC_URL, default (local API).
func LobbyIssuer() string {
	if v := normalizeIssuerURL(os.Getenv("LOBBY_ISSUER_URL")); v != "" {
		return v
	}
	if v := normalizeIssuerURL(os.Getenv("LOBBY_PUBLIC_URL")); v != "" {
		return v
	}
	return defaultLobbyIssuer
}

// LobbyReturnURL is where games send players after a match (browser-facing Lobby URL).
// Resolution order: LOBBY_RETURN_URL, LOBBY_PUBLIC_URL, LobbyIssuer().
func LobbyReturnURL() string {
	if v := normalizeIssuerURL(os.Getenv("LOBBY_RETURN_URL")); v != "" {
		return v
	}
	if v := normalizeIssuerURL(os.Getenv("LOBBY_PUBLIC_URL")); v != "" {
		return v
	}
	return LobbyIssuer()
}

// LobbyGraphQLURL is the Lobby GraphQL endpoint (player lookup, match reporting, etc.).
func LobbyGraphQLURL() string {
	return LobbyIssuer() + "/graphql"
}

func normalizeIssuerURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return strings.TrimRight(raw, "/")
	}
	u.Fragment = ""
	u.RawQuery = ""
	u.Path = strings.TrimRight(u.Path, "/")
	out := u.String()
	return strings.TrimRight(out, "/")
}
