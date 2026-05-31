package store

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestPickDistinctWaitingEntriesPreventsSelfMatch(t *testing.T) {
	userID := uuid.New()
	entries := []QueueEntry{
		{ID: uuid.New(), UserID: userID, Status: QueueStatusWaiting},
		{ID: uuid.New(), UserID: userID, Status: QueueStatusWaiting},
		{ID: uuid.New(), UserID: uuid.New(), Status: QueueStatusWaiting},
	}
	picked := pickDistinctWaitingEntries(entries, 2)
	if len(picked) != 2 {
		t.Fatalf("expected 2 distinct picks, got %d", len(picked))
	}
	if picked[0].UserID == picked[1].UserID {
		t.Fatal("expected distinct user ids in match roster")
	}
}

func TestJoinModeQueueIdempotentWhileWaiting(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	queueID := DemoDefaultQueueID
	user, err := st.CreateUser(ctx, CreateUserParams{Email: "idem-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(user.ID)

	first, err := st.JoinModeQueue(ctx, queueID, user.ID)
	if err != nil {
		t.Fatalf("first join: %v", err)
	}
	if first.Status != QueueStatusWaiting {
		t.Fatalf("expected waiting, got %s", first.Status)
	}

	second, err := st.JoinModeQueue(ctx, queueID, user.ID)
	if err != nil {
		t.Fatalf("second join: %v", err)
	}
	if !second.AlreadyInQueue {
		t.Fatalf("expected AlreadyInQueue on re-join")
	}
}

func TestJoinModeQueueRejectsWhenAlreadyMatched(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	queueID := DemoDefaultQueueID

	userA, err := st.CreateUser(ctx, CreateUserParams{Email: "matched-a-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser A: %v", err)
	}
	cleaner.TrackUser(userA.ID)

	userB, err := st.CreateUser(ctx, CreateUserParams{Email: "matched-b-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser B: %v", err)
	}
	cleaner.TrackUser(userB.ID)

	if _, err := st.JoinModeQueue(ctx, queueID, userA.ID); err != nil {
		t.Fatalf("join A: %v", err)
	}
	if _, err := st.JoinModeQueue(ctx, queueID, userB.ID); err != nil {
		t.Fatalf("join B: %v", err)
	}

	_, err = st.JoinModeQueue(ctx, queueID, userA.ID)
	if err == nil {
		t.Fatal("expected error when joining queue while already matched")
	}
}
