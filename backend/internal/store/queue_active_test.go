package store

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestGetUserActiveIntentPrefersActiveSessionOverWaiting(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	queueID := DemoDefaultQueueID
	user, err := st.CreateUser(ctx, CreateUserParams{Email: "active-over-wait-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(user.ID)

	if _, err := st.JoinModeQueue(ctx, queueID, user.ID, "", nil); err != nil {
		t.Fatalf("JoinModeQueue waiting: %v", err)
	}

	modes, err := st.ListGameModesByGameID(ctx, DemoPrimaryGameID)
	if err != nil {
		t.Fatalf("ListGameModesByGameID: %v", err)
	}
	if len(modes) == 0 {
		t.Fatal("expected demo game modes")
	}
	modeID := modes[0].ID

	var sessionID uuid.UUID
	err = st.db.QueryRowContext(ctx, `
		INSERT INTO game_sessions (game_id, status, mode_id, mode_queue_id, started_at)
		VALUES ($1, 'active', $2, $3, NOW())
		RETURNING id
	`, DemoPrimaryGameID, modeID, queueID).Scan(&sessionID)
	if err != nil {
		t.Fatalf("insert active session: %v", err)
	}
	if err := st.AddSessionParticipant(ctx, sessionID, user.ID, "player1"); err != nil {
		t.Fatalf("AddSessionParticipant: %v", err)
	}

	active, err := st.GetUserActiveIntent(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetUserActiveIntent: %v", err)
	}
	if active == nil || !active.Matched || active.Waiting {
		t.Fatalf("expected matched active session, got %+v", active)
	}
	if active.SessionID == nil || *active.SessionID != sessionID {
		t.Fatalf("session = %v, want %s", active.SessionID, sessionID)
	}
}
