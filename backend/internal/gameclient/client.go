package gameclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// AssignmentSeat is one seat in a Lobby-pushed roster.
type AssignmentSeat struct {
	SeatKey     string `json:"seatKey"`
	LobbyUserID string `json:"lobbyUserId"`
	DisplayName string `json:"displayName,omitempty"`
	Team        string `json:"team,omitempty"`
	Role        string `json:"role,omitempty"`
}

// Assignment is the POST /api/v1/matches body (option 2 push provisioning).
type Assignment struct {
	ExternalMatchID string           `json:"externalMatchId"`
	GameMode        string           `json:"gameMode"`
	BestOf          int              `json:"bestOf,omitempty"`
	Seats           []AssignmentSeat `json:"seats"`
}

// BannedPlayersError is returned when the game rejects the roster (HTTP 403).
type BannedPlayersError struct {
	BannedLobbyUserIDs []string
	Message            string
}

func (e *BannedPlayersError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("roster rejected: banned %v", e.BannedLobbyUserIDs)
}

// Client pushes match assignments to a third-party game API.
type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// ProvisionMatch creates or confirms a match on the game server (idempotent on externalMatchId).
func (c *Client) ProvisionMatch(ctx context.Context, apiBaseURL, serviceToken string, assignment Assignment) error {
	base := strings.TrimRight(strings.TrimSpace(apiBaseURL), "/")
	if base == "" {
		return errors.New("gameclient: api base URL is required")
	}

	body, err := json.Marshal(map[string]any{"assignment": assignment})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/api/v1/matches", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if header := serviceAuthHeader(serviceToken); header != "" {
		req.Header.Set("Authorization", header)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("gameclient: provision match: %w", err)
	}
	defer resp.Body.Close()

	payload, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if resp.StatusCode == http.StatusForbidden {
		var errBody struct {
			Error              string   `json:"error"`
			BannedLobbyUserIDs []string `json:"bannedLobbyUserIds"`
		}
		_ = json.Unmarshal(payload, &errBody)
		return &BannedPlayersError{
			BannedLobbyUserIDs: errBody.BannedLobbyUserIDs,
			Message:            errBody.Error,
		}
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("gameclient: provision match: HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(payload)))
	}
	return nil
}

func serviceAuthHeader(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return ""
	}
	if len(token) >= 7 && strings.EqualFold(token[:7], "bearer ") {
		return token
	}
	return "Bearer " + token
}
