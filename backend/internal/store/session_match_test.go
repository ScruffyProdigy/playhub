package store

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRequeueAfterReturnUsesNewSessionNotStaleActive(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	queueID := DemoDefaultQueueID
	userA, err := st.CreateUser(ctx, CreateUserParams{Email: "requeue-a-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser A: %v", err)
	}
	cleaner.TrackUser(userA.ID)
	userB, err := st.CreateUser(ctx, CreateUserParams{Email: "requeue-b-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser B: %v", err)
	}
	cleaner.TrackUser(userB.ID)

	if _, err := st.JoinModeQueue(ctx, queueID, userA.ID); err != nil {
		t.Fatalf("first join A: %v", err)
	}
	if _, err := st.JoinModeQueue(ctx, queueID, userB.ID); err != nil {
		t.Fatalf("first join B: %v", err)
	}

	first, err := st.GetUserActiveQueue(ctx, userA.ID)
	if err != nil {
		t.Fatalf("GetUserActiveQueue after first match: %v", err)
	}
	if first == nil || first.SessionID == nil {
		t.Fatal("expected first matched session")
	}
	firstSessionID := *first.SessionID

	// Simulate return hub clearing matched rows without reportMatchResult completing the session.
	if err := st.ReleaseUserMatchedQueue(ctx, queueID, userA.ID); err != nil {
		t.Fatalf("ReleaseUserMatchedQueue A: %v", err)
	}
	if err := st.ReleaseUserMatchedQueue(ctx, queueID, userB.ID); err != nil {
		t.Fatalf("ReleaseUserMatchedQueue B: %v", err)
	}

	firstSession, err := st.GetSessionByID(ctx, firstSessionID)
	if err != nil {
		t.Fatalf("GetSessionByID first: %v", err)
	}
	if firstSession.Status != "active" {
		t.Fatalf("first session should still be active before re-queue, got %q", firstSession.Status)
	}

	if _, err := st.JoinModeQueue(ctx, queueID, userA.ID); err != nil {
		t.Fatalf("second join A: %v", err)
	}
	second, err := st.JoinModeQueue(ctx, queueID, userB.ID)
	if err != nil {
		t.Fatalf("second join B: %v", err)
	}
	if second.SessionID == nil {
		t.Fatal("expected session on second match")
	}
	if *second.SessionID == firstSessionID {
		t.Fatalf("re-queue reused stale session %s", firstSessionID)
	}

	active, err := st.GetUserActiveQueue(ctx, userA.ID)
	if err != nil {
		t.Fatalf("GetUserActiveQueue after re-queue: %v", err)
	}
	if active == nil || active.SessionID == nil || *active.SessionID != *second.SessionID {
		t.Fatalf("active queue session = %v, want %s", active, second.SessionID)
	}

	firstSession, err = st.GetSessionByID(ctx, firstSessionID)
	if err != nil {
		t.Fatalf("GetSessionByID first after re-queue: %v", err)
	}
	if firstSession.Status != "completed" {
		t.Fatalf("stale session status = %q, want completed", firstSession.Status)
	}
}

func TestGetMatchedSessionForUserAndModeQueuePicksNewestSession(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	queueID := DemoDefaultQueueID
	user, err := st.CreateUser(ctx, CreateUserParams{Email: "newest-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(user.ID)

	modes, err := st.ListGameModesByGameID(ctx, DemoRPSGameID)
	if err != nil {
		t.Fatalf("ListGameModesByGameID: %v", err)
	}
	if len(modes) == 0 {
		t.Fatal("expected demo game modes")
	}
	modeID := modes[0].ID

	insertSession := func(startedAt time.Time) uuid.UUID {
		t.Helper()
		var id uuid.UUID
		err := st.db.QueryRowContext(ctx, `
			INSERT INTO game_sessions (game_id, status, mode_id, mode_queue_id, started_at)
			VALUES ($1, 'active', $2, $3, $4)
			RETURNING id
		`, DemoRPSGameID, modeID, queueID, startedAt).Scan(&id)
		if err != nil {
			t.Fatalf("insert session: %v", err)
		}
		if err := st.AddSessionParticipant(ctx, id, user.ID, "player1"); err != nil {
			t.Fatalf("AddSessionParticipant: %v", err)
		}
		return id
	}

	oldSessionID := insertSession(time.Now().Add(-time.Hour))
	newSessionID := insertSession(time.Now())

	_, err = st.db.ExecContext(ctx, `
		INSERT INTO game_queues (game_id, user_id, mode_queue_id, status, matched_at)
		VALUES ($1, $2, $3, 'matched', NOW())
	`, DemoRPSGameID, user.ID, queueID)
	if err != nil {
		t.Fatalf("insert matched queue row: %v", err)
	}

	got, err := st.GetMatchedSessionForUserAndModeQueue(ctx, queueID, user.ID)
	if err != nil {
		t.Fatalf("GetMatchedSessionForUserAndModeQueue: %v", err)
	}
	if got.ID != newSessionID {
		t.Fatalf("session = %s, want newest %s (old=%s)", got.ID, newSessionID, oldSessionID)
	}
}
