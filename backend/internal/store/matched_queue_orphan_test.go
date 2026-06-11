package store

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestJoinModeQueueHealsStaleMatchedRowWithoutSession(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	queueID := DemoDefaultQueueID
	user, err := st.CreateUser(ctx, CreateUserParams{Email: "stale-matched-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(user.ID)

	if _, err := st.JoinModeQueue(ctx, queueID, user.ID, "", nil); err != nil {
		t.Fatalf("join: %v", err)
	}

	var queueRowID uuid.UUID
	if err := st.db.QueryRowContext(ctx, `
		UPDATE game_queues
		SET status = 'matched', matched_at = NOW()
		WHERE user_id = $1 AND mode_queue_id = $2 AND status = 'waiting'
		RETURNING id
	`, user.ID, queueID).Scan(&queueRowID); err != nil {
		t.Fatalf("simulate stale matched row: %v", err)
	}

	intent, err := st.GetUserActiveIntent(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetUserActiveIntent: %v", err)
	}
	if intent != nil && intent.Matched {
		t.Fatalf("expected no matched intent for orphan row, got %+v", intent)
	}

	if _, err := st.JoinModeQueue(ctx, queueID, user.ID, "", nil); err != nil {
		t.Fatalf("re-join after stale matched heal: %v", err)
	}

	var status string
	if err := st.db.QueryRowContext(ctx, `
		SELECT status FROM game_queues WHERE id = $1
	`, queueRowID).Scan(&status); err != nil {
		t.Fatalf("old row status: %v", err)
	}
	if status != "cancelled" {
		t.Fatalf("stale matched row status = %q, want cancelled", status)
	}
}

func TestGetUserActiveIntentPrefersActiveSessionOverStaleMatchedRow(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	queueID := DemoDefaultQueueID
	userA, err := st.CreateUser(ctx, CreateUserParams{Email: "live-match-a-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser A: %v", err)
	}
	cleaner.TrackUser(userA.ID)
	userB, err := st.CreateUser(ctx, CreateUserParams{Email: "live-match-b-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser B: %v", err)
	}
	cleaner.TrackUser(userB.ID)

	if _, err := st.JoinModeQueue(ctx, queueID, userA.ID, "", nil); err != nil {
		t.Fatalf("join A: %v", err)
	}
	if _, err := st.JoinModeQueue(ctx, queueID, userB.ID, "", nil); err != nil {
		t.Fatalf("join B: %v", err)
	}
	mustReconcileForming(t, st, ctx, queueID)

	intent, err := st.GetUserActiveIntent(ctx, userA.ID)
	if err != nil {
		t.Fatalf("GetUserActiveIntent: %v", err)
	}
	if intent == nil || !intent.Matched || intent.SessionID == nil {
		t.Fatal("expected live matched intent with session")
	}

	if _, err := st.JoinModeQueue(ctx, queueID, userA.ID, "", nil); err == nil {
		t.Fatal("expected join to fail while session is active")
	}
}
