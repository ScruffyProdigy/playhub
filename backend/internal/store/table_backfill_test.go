package store

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// Two kings in different rooms each sit local seat "1" and start backfill; they should match.
func TestStartTableBackfillTwoPartialDuelsMatch(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	game, mode := setupDuelMode(t, st, cleaner)
	queues, err := st.ListModeQueuesByModeID(ctx, mode.ID)
	if err != nil {
		t.Fatalf("ListModeQueuesByModeID: %v", err)
	}
	queueID := queues[0].ID

	userA, err := st.CreateUser(ctx, CreateUserParams{Email: "backfill-a-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser A: %v", err)
	}
	cleaner.TrackUser(userA.ID)
	userB, err := st.CreateUser(ctx, CreateUserParams{Email: "backfill-b-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser B: %v", err)
	}
	cleaner.TrackUser(userB.ID)

	roomA, err := st.CreateRoom(ctx, userA.ID)
	if err != nil {
		t.Fatalf("CreateRoom A: %v", err)
	}
	roomB, err := st.CreateRoom(ctx, userB.ID)
	if err != nil {
		t.Fatalf("CreateRoom B: %v", err)
	}

	tableA, err := st.CreateTable(ctx, roomA.ID, game.ID, mode.ID, userA.ID)
	if err != nil {
		t.Fatalf("CreateTable A: %v", err)
	}
	tableB, err := st.CreateTable(ctx, roomB.ID, game.ID, mode.ID, userB.ID)
	if err != nil {
		t.Fatalf("CreateTable B: %v", err)
	}

	seats, err := st.ListGameModeSeats(ctx, mode.ID)
	if err != nil {
		t.Fatalf("ListGameModeSeats: %v", err)
	}
	firstSeat := seats[0].SeatKey

	if _, err := st.SitAtTable(ctx, tableA.ID, userA.ID, firstSeat); err != nil {
		t.Fatalf("sit A: %v", err)
	}
	if _, err := st.SitAtTable(ctx, tableB.ID, userB.ID, firstSeat); err != nil {
		t.Fatalf("sit B: %v", err)
	}

	if _, err := st.StartTableBackfill(ctx, tableA.ID, userA.ID, queueID); err != nil {
		t.Fatalf("backfill A: %v", err)
	}
	recA := mustReconcileForming(t, st, ctx, queueID)
	if recA.Fired {
		t.Fatalf("expected A waiting, got fired session %s", recA.SessionID)
	}

	if _, err := st.StartTableBackfill(ctx, tableB.ID, userB.ID, queueID); err != nil {
		t.Fatalf("backfill B: %v", err)
	}
	rec := mustReconcileForming(t, st, ctx, queueID)
	if !rec.Fired || rec.SessionID == nil {
		t.Fatalf("expected B to complete match, got %+v", rec)
	}

	participants, err := st.ListSessionSeatAssignments(ctx, *rec.SessionID)
	if err != nil {
		t.Fatalf("ListSessionSeatAssignments: %v", err)
	}
	if len(participants) != 2 {
		t.Fatalf("expected 2 participants, got %d", len(participants))
	}
}
