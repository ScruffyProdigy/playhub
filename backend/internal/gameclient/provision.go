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

// MatchProvisioner pushes roster assignments to game servers.
type MatchProvisioner interface {
	ProvisionMatch(ctx context.Context, req ProvisionRequest) error
}
