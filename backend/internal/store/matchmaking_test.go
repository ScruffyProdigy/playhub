package store

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestJoinGameQueueMatchesWhenMinPlayersReached(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	gameID := uuid.MustParse("a1000000-0000-4000-8000-000000000001")

	userA, err := st.CreateUser(ctx, CreateUserParams{Email: "match-a-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser A failed: %v", err)
	}
	cleaner.TrackUser(userA.ID)

	userB, err := st.CreateUser(ctx, CreateUserParams{Email: "match-b-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser B failed: %v", err)
	}
	cleaner.TrackUser(userB.ID)

	first, err := st.JoinGameQueue(ctx, gameID, userA.ID)
	if err != nil {
		t.Fatalf("JoinGameQueue A failed: %v", err)
	}
	if first.Status != QueueStatusWaiting {
		t.Fatalf("expected A to be waiting, got %s", first.Status)
	}

	second, err := st.JoinGameQueue(ctx, gameID, userB.ID)
	if err != nil {
		t.Fatalf("JoinGameQueue B failed: %v", err)
	}
	if second.Status != QueueStatusMatched || second.SessionID == nil {
		t.Fatalf("expected B to be matched, got %+v", second)
	}

	if len(second.NotifyUserIDs) != 2 {
		t.Fatalf("expected 2 users notified, got %d", len(second.NotifyUserIDs))
	}

	view, err := st.GetUserQueueView(ctx, gameID, userA.ID)
	if err != nil {
		t.Fatalf("GetUserQueueView A failed: %v", err)
	}
	if !view.Matched || view.SessionID == nil {
		t.Fatalf("expected A to see existing match, got %+v", view)
	}
	if *view.SessionID != *second.SessionID {
		t.Fatalf("expected same session id, got %s vs %s", view.SessionID, second.SessionID)
	}
}

func TestJoinGameQueueDoesNotReuseStaleActiveSession(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	gameID := uuid.MustParse("a1000000-0000-4000-8000-000000000001")

	user, err := st.CreateUser(ctx, CreateUserParams{Email: "stale-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	cleaner.TrackUser(user.ID)

	session, err := st.CreateSession(ctx, gameID)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}
	if err := st.AddSessionParticipant(ctx, session.ID, user.ID, "player"); err != nil {
		t.Fatalf("AddSessionParticipant failed: %v", err)
	}

	result, err := st.JoinGameQueue(ctx, gameID, user.ID)
	if err != nil {
		t.Fatalf("JoinGameQueue failed: %v", err)
	}
	if result.Status != QueueStatusWaiting {
		t.Fatalf("expected waiting for new queue join, got %s", result.Status)
	}
	if result.SessionID != nil {
		t.Fatalf("expected no session id, got %v", result.SessionID)
	}
}

func TestJoinGameQueueIdempotentWhileWaiting(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	gameID := uuid.MustParse("a1000000-0000-4000-8000-000000000001")
	user, err := st.CreateUser(ctx, CreateUserParams{Email: "idem-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(user.ID)

	first, err := st.JoinGameQueue(ctx, gameID, user.ID)
	if err != nil {
		t.Fatalf("first join: %v", err)
	}
	if first.Status != QueueStatusWaiting {
		t.Fatalf("expected waiting, got %s", first.Status)
	}

	second, err := st.JoinGameQueue(ctx, gameID, user.ID)
	if err != nil {
		t.Fatalf("second join: %v", err)
	}
	if !second.AlreadyInQueue {
		t.Fatalf("expected AlreadyInQueue on re-join")
	}
	if second.Status != QueueStatusWaiting {
		t.Fatalf("expected still waiting, got %s", second.Status)
	}
}

func TestJoinGameQueueRejectsWhenAlreadyMatched(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	gameID := uuid.MustParse("a1000000-0000-4000-8000-000000000001")

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

	if _, err := st.JoinGameQueue(ctx, gameID, userA.ID); err != nil {
		t.Fatalf("join A: %v", err)
	}
	if _, err := st.JoinGameQueue(ctx, gameID, userB.ID); err != nil {
		t.Fatalf("join B: %v", err)
	}

	_, err = st.JoinGameQueue(ctx, gameID, userA.ID)
	if err == nil {
		t.Fatal("expected error when joining queue while already matched")
	}
	if !errors.Is(err, ErrAlreadyMatched) {
		t.Fatalf("expected ErrAlreadyMatched, got %v", err)
	}
}

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
