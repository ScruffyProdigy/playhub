package store

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCompleteSessionReleasesMatchedQueueRows(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	queueID := DemoDefaultQueueID
	userA, err := st.CreateUser(ctx, CreateUserParams{Email: "complete-a-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser A: %v", err)
	}
	cleaner.TrackUser(userA.ID)
	userB, err := st.CreateUser(ctx, CreateUserParams{Email: "complete-b-" + uuid.NewString() + "@example.com"})
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

	view, err := st.GetUserActiveQueue(ctx, userA.ID)
	if err != nil {
		t.Fatalf("GetUserActiveQueue: %v", err)
	}
	if view == nil || !view.Matched || view.SessionID == nil {
		t.Fatal("expected matched session")
	}
	sessionID := *view.SessionID

	if err := st.CompleteSession(ctx, sessionID, time.Now()); err != nil {
		t.Fatalf("CompleteSession: %v", err)
	}

	session, err := st.GetSessionByID(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetSessionByID: %v", err)
	}
	if session.Status != "completed" {
		t.Fatalf("session status = %q, want completed", session.Status)
	}

	for _, uid := range []uuid.UUID{userA.ID, userB.ID} {
		active, err := st.GetUserActiveQueue(ctx, uid)
		if err != nil {
			t.Fatalf("GetUserActiveQueue(%s): %v", uid, err)
		}
		if active != nil && active.Matched {
			t.Fatalf("user %s still matched after CompleteSession", uid)
		}
	}
}
