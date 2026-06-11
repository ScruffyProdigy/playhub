package gameclient

import "context"

// ProvisionRequest is the server-to-server match push to a game API.
type ProvisionRequest struct {
	APIBaseURL   string
	ServiceToken string
	LobbyID      string
	Lobby        LobbyInfo
	Assignment   Assignment
}

// ProvisionResult is parsed from a successful game provision response.
type ProvisionResult struct {
	// LaunchURLs maps lobbyUserId → game-minted URL base (no JWT).
	LaunchURLs map[string]string
	// LaunchURLTemplate is an optional per-seat template with placeholders.
	LaunchURLTemplate string
}

// MatchProvisioner pushes roster assignments to game servers.
type MatchProvisioner interface {
	ProvisionMatch(ctx context.Context, req ProvisionRequest) (ProvisionResult, error)
}
