package store

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestJoinModeQueueHealsOrphanedPartyWithoutWaitingRow(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	queueID := DemoDefaultQueueID
	user, err := st.CreateUser(ctx, CreateUserParams{Email: "orphan-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(user.ID)

	first, err := st.JoinModeQueue(ctx, queueID, user.ID, "", nil)
	if err != nil {
		t.Fatalf("first join: %v", err)
	}
	if first.Status != QueueStatusWaiting {
		t.Fatalf("expected waiting, got %s", first.Status)
	}

	var partyID uuid.UUID
	if err := st.db.QueryRowContext(ctx, `
		SELECT party_id FROM game_queues
		WHERE user_id = $1 AND mode_queue_id = $2 AND status = 'waiting'
	`, user.ID, queueID).Scan(&partyID); err != nil {
		t.Fatalf("party_id: %v", err)
	}
	if partyID == uuid.Nil {
		t.Fatal("expected party on waiting queue row")
	}

	if _, err := st.db.ExecContext(ctx, `
		UPDATE game_queues SET status = 'cancelled'
		WHERE user_id = $1 AND mode_queue_id = $2 AND status = 'waiting'
	`, user.ID, queueID); err != nil {
		t.Fatalf("simulate orphan queue cancel: %v", err)
	}

	second, err := st.JoinModeQueue(ctx, queueID, user.ID, "", nil)
	if err != nil {
		t.Fatalf("re-join after orphan: %v", err)
	}
	if second.Status != QueueStatusWaiting {
		t.Fatalf("expected waiting after heal, got %s", second.Status)
	}

	var partyStatus string
	if err := st.db.QueryRowContext(ctx, `SELECT status FROM parties WHERE id = $1`, partyID).Scan(&partyStatus); err != nil {
		t.Fatalf("old party status: %v", err)
	}
	if partyStatus != PartyStatusCancelled {
		t.Fatalf("old party status = %q, want cancelled", partyStatus)
	}
}

func TestLeaveModeQueueCancelsActiveParty(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	queueID := DemoDefaultQueueID
	user, err := st.CreateUser(ctx, CreateUserParams{Email: "leave-party-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(user.ID)

	if _, err := st.JoinModeQueue(ctx, queueID, user.ID, "", nil); err != nil {
		t.Fatalf("join: %v", err)
	}

	var partyID uuid.UUID
	if err := st.db.QueryRowContext(ctx, `
		SELECT party_id FROM game_queues
		WHERE user_id = $1 AND mode_queue_id = $2 AND status = 'waiting'
	`, user.ID, queueID).Scan(&partyID); err != nil {
		t.Fatalf("party_id: %v", err)
	}

	if _, err := st.LeaveModeQueue(ctx, queueID, user.ID); err != nil {
		t.Fatalf("LeaveModeQueue: %v", err)
	}

	var partyStatus string
	if err := st.db.QueryRowContext(ctx, `SELECT status FROM parties WHERE id = $1`, partyID).Scan(&partyStatus); err != nil {
		t.Fatalf("party status: %v", err)
	}
	if partyStatus != PartyStatusCancelled {
		t.Fatalf("party status = %q, want cancelled", partyStatus)
	}

	third, err := st.JoinModeQueue(ctx, queueID, user.ID, "", nil)
	if err != nil {
		t.Fatalf("re-join after leave: %v", err)
	}
	if third.Status != QueueStatusWaiting {
		t.Fatalf("expected waiting, got %s", third.Status)
	}
}

func TestJoinModeQueueRejectsAlreadyInActivePartyMessage(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	queueID := DemoDefaultQueueID
	user, err := st.CreateUser(ctx, CreateUserParams{Email: "party-msg-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(user.ID)

	if _, err := st.JoinModeQueue(ctx, queueID, user.ID, "", nil); err != nil {
		t.Fatalf("join: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, `
		UPDATE game_queues SET status = 'cancelled'
		WHERE user_id = $1 AND status = 'waiting'
	`, user.ID); err != nil {
		t.Fatalf("orphan queue: %v", err)
	}

	_, err = st.JoinModeQueue(ctx, queueID, user.ID, "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "already in an active party") {
			t.Fatalf("expected orphan heal, got stale party block: %v", err)
		}
		t.Fatalf("re-join: %v", err)
	}
}
