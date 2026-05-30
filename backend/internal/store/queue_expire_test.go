package store

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestExpireStaleMatchedQueueClearsOldMatch(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	gameID := uuid.MustParse("a1000000-0000-4000-8000-000000000001")
	userA, err := st.CreateUser(ctx, CreateUserParams{Email: "stale-exp-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(userA.ID)
	userB, err := st.CreateUser(ctx, CreateUserParams{Email: "stale-exp-b-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser B: %v", err)
	}
	cleaner.TrackUser(userB.ID)

	if _, err := st.JoinGameQueue(ctx, gameID, userA.ID); err != nil {
		t.Fatalf("join A: %v", err)
	}
	if _, err := st.JoinGameQueue(ctx, gameID, userB.ID); err != nil {
		t.Fatalf("join B: %v", err)
	}

	_, err = st.db.ExecContext(ctx, `
		UPDATE game_queues
		SET matched_at = NOW() - INTERVAL '2 hours'
		WHERE game_id = $1 AND user_id = $2 AND status = 'matched'
	`, gameID, userA.ID)
	if err != nil {
		t.Fatalf("backdate matched_at: %v", err)
	}

	view, err := st.GetUserQueueView(ctx, gameID, userA.ID)
	if err != nil {
		t.Fatalf("GetUserQueueView: %v", err)
	}
	if view.Matched || view.Waiting {
		t.Fatalf("expected stale match cleared, got %+v", view)
	}
}
