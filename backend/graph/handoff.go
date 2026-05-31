package graph

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/auth"
	"github.com/scruffyprodigy/playhub/internal/gameclient"
	"github.com/scruffyprodigy/playhub/internal/store"
)

func (r *Resolver) gameProvisioner() gameclient.MatchProvisioner {
	if r.GameProvisioner != nil {
		return r.GameProvisioner
	}
	return gameclient.NewClient()
}

func handoffDebugEnabled() bool {
	v := strings.TrimSpace(os.Getenv("LOBBY_HANDOFF_DEBUG"))
	return v == "1" || strings.EqualFold(v, "true")
}

func gamePlayURL(game *store.Game, fallback string) string {
	if game != nil && game.PlayURL != nil {
		if u := strings.TrimSpace(*game.PlayURL); u != "" {
			return strings.TrimRight(u, "/")
		}
	}
	return strings.TrimRight(strings.TrimSpace(fallback), "/")
}

// resolvedAPIBaseURL prefers the catalog row, then GAME_API_BASE_URL (local dev fallback).
func (r *Resolver) resolvedAPIBaseURL(game *store.Game) string {
	if game != nil && game.APIBaseURL != nil {
		if u := strings.TrimSpace(*game.APIBaseURL); u != "" {
			return strings.TrimRight(u, "/")
		}
	}
	return strings.TrimRight(strings.TrimSpace(os.Getenv("GAME_API_BASE_URL")), "/")
}

func modeKeyFromCatalog(mode *store.GameMode) (string, error) {
	if mode == nil || strings.TrimSpace(mode.ModeKey) == "" {
		return "", fmt.Errorf("catalog mode is required")
	}
	return strings.TrimSpace(mode.ModeKey), nil
}

func assignmentFromParticipants(sessionID uuid.UUID, mode *store.GameMode, participants []store.SessionParticipant) (gameclient.Assignment, error) {
	modeKey, err := modeKeyFromCatalog(mode)
	if err != nil {
		return gameclient.Assignment{}, err
	}
	assignment := gameclient.Assignment{
		ExternalMatchID: sessionID.String(),
		GameMode:        modeKey,
		Seats:           make([]gameclient.AssignmentSeat, 0, len(participants)),
	}
	for _, p := range participants {
		assignment.Seats = append(assignment.Seats, gameclient.AssignmentSeat{
			SeatKey:     p.SeatKey,
			LobbyUserID: p.UserID.String(),
		})
	}
	return assignment, nil
}

func lobbyProvisionInfo() gameclient.LobbyInfo {
	return gameclient.LobbyInfo{
		ReturnURL:    auth.LobbyReturnURL(),
		GraphqlURL:   auth.LobbyGraphQLURL(),
		ServiceToken: auth.GameServiceTokenFromEnv(),
	}
}

func (r *Resolver) provisionParticipantsOnGame(ctx context.Context, game *store.Game, sessionID uuid.UUID, participants []store.SessionParticipant) error {
	apiBase := r.resolvedAPIBaseURL(game)
	if apiBase == "" {
		return fmt.Errorf("game API base URL is not configured (set games.api_base_url or GAME_API_BASE_URL)")
	}
	if len(participants) == 0 {
		return fmt.Errorf("cannot provision match with no seated players")
	}

	st, err := r.requireStore()
	if err != nil {
		return err
	}

	mode, err := st.GetGameModeForSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("catalog mode for session: %w", err)
	}
	assignment, err := assignmentFromParticipants(sessionID, mode, participants)
	if err != nil {
		return err
	}

	if handoffDebugEnabled() {
		log.Printf("handoff: POST %s/api/v1/matches externalMatchId=%s seats=%d", apiBase, sessionID, len(assignment.Seats))
	}

	lobby := lobbyProvisionInfo()
	if err := r.gameProvisioner().ProvisionMatch(ctx, gameclient.ProvisionRequest{
		APIBaseURL:   apiBase,
		ServiceToken: lobby.ServiceToken,
		LobbyID:      auth.LobbyIssuer(),
		Lobby:        lobby,
		Assignment:   assignment,
	}); err != nil {
		if banned, ok := err.(*gameclient.BannedPlayersError); ok {
			st, stErr := r.requireStore()
			if stErr != nil {
				return stErr
			}
			bannedIDs := parseBannedLobbyUserIDs(banned.BannedLobbyUserIDs)
			if rbErr := st.RollbackMatchedSession(ctx, sessionID, bannedIDs); rbErr != nil {
				return fmt.Errorf("provision banned and rollback failed: %w (rollback: %v)", err, rbErr)
			}
		}
		return err
	}
	return nil
}

// finalizeMatchedSession pushes the roster to the game and returns per-user launch URLs.
func (r *Resolver) finalizeMatchedSession(ctx context.Context, game *store.Game, sessionID uuid.UUID, notifyUserIDs []uuid.UUID) (map[uuid.UUID]string, error) {
	st, err := r.requireStore()
	if err != nil {
		return nil, err
	}
	authService, err := r.requireAuth()
	if err != nil {
		return nil, err
	}

	participants, err := st.ListSessionSeatAssignments(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	playURL := gamePlayURL(game, r.GameClientBaseURL)
	if playURL == "" {
		return nil, fmt.Errorf("game play URL is not configured")
	}
	audience := r.resolvedAPIBaseURL(game)
	if audience == "" {
		return nil, fmt.Errorf("game API base URL is not configured")
	}

	if err := r.provisionParticipantsOnGame(ctx, game, sessionID, participants); err != nil {
		return nil, err
	}

	signer := authService.Signer()
	externalMatchID := sessionID.String()
	urls := make(map[uuid.UUID]string, len(participants))
	for _, p := range participants {
		launch, err := launchURLForSeat(signer, playURL, audience, externalMatchID, p.UserID, p.SeatKey, p.DisplayName)
		if err != nil {
			return nil, err
		}
		urls[p.UserID] = launch
	}

	for _, uid := range notifyUserIDs {
		if _, ok := urls[uid]; !ok {
			return nil, fmt.Errorf("missing launch url for notified user %s", uid)
		}
	}
	return urls, nil
}

func launchURLForSeat(signer *auth.Signer, playURL, audience, externalMatchID string, userID uuid.UUID, seatKey, displayName string) (string, error) {
	token, err := signer.SignSeatToken(userID, audience, externalMatchID, seatKey, displayName, 0)
	if err != nil {
		return "", err
	}
	return buildLaunchURL(playURL, externalMatchID, token)
}

func buildLaunchURL(playURL, externalMatchID, token string) (string, error) {
	u, err := url.Parse(playURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("match", externalMatchID)
	q.Set("token", token)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func parseBannedLobbyUserIDs(ids []string) []uuid.UUID {
	out := make([]uuid.UUID, 0, len(ids))
	for _, s := range ids {
		id, err := uuid.Parse(strings.TrimSpace(s))
		if err != nil {
			continue
		}
		out = append(out, id)
	}
	return out
}

// signLaunchURL mints a fresh seat JWT and launch link. The match must already
// have been provisioned on the game (finalizeMatchedSession); this does not POST again.
func (r *Resolver) signLaunchURL(ctx context.Context, game *store.Game, sessionID, userID uuid.UUID) (string, error) {
	st, err := r.requireStore()
	if err != nil {
		return "", err
	}
	authService, err := r.requireAuth()
	if err != nil {
		return "", err
	}

	participants, err := st.ListSessionSeatAssignments(ctx, sessionID)
	if err != nil {
		return "", err
	}

	playURL := gamePlayURL(game, r.GameClientBaseURL)
	if playURL == "" {
		return "", fmt.Errorf("game play URL is not configured")
	}
	audience := r.resolvedAPIBaseURL(game)
	if audience == "" {
		return "", fmt.Errorf("game API base URL is not configured")
	}

	for _, p := range participants {
		if p.UserID == userID {
			return launchURLForSeat(authService.Signer(), playURL, audience, sessionID.String(), userID, p.SeatKey, p.DisplayName)
		}
	}
	return "", fmt.Errorf("user is not seated in session")
}
