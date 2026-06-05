package store

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/gameclient"
)

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

	first, err := st.JoinModeQueue(ctx, queueID, user.ID, "")
	if err != nil {
		t.Fatalf("first join: %v", err)
	}
	if first.Status != QueueStatusWaiting {
		t.Fatalf("expected waiting, got %s", first.Status)
	}

	second, err := st.JoinModeQueue(ctx, queueID, user.ID, "")
	if err != nil {
		t.Fatalf("second join: %v", err)
	}
	if !second.AlreadyInQueue {
		t.Fatalf("expected AlreadyInQueue on re-join")
	}
}

func TestJoinModeQueueReplacesStaleWaitingRowForSameGame(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	queueID := DemoDefaultQueueID
	user, err := st.CreateUser(ctx, CreateUserParams{Email: "stale-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(user.ID)

	// Simulate a legacy waiting row (no mode_queue_id) that blocks the per-game unique index.
	_, err = st.db.ExecContext(ctx, `
		INSERT INTO game_queues (game_id, user_id, status)
		VALUES ($1, $2, 'waiting')
	`, DemoRPSGameID, user.ID)
	if err != nil {
		t.Fatalf("insert legacy queue row: %v", err)
	}
	result, err := st.JoinModeQueue(ctx, queueID, user.ID, "")
	if err != nil {
		t.Fatalf("JoinModeQueue after legacy row: %v", err)
	}
	if result.Status != QueueStatusWaiting {
		t.Fatalf("expected waiting, got %s", result.Status)
	}

	var waitingCount int
	if err := st.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM game_queues
		WHERE user_id = $1 AND status = 'waiting' AND mode_queue_id = $2
	`, user.ID, queueID).Scan(&waitingCount); err != nil {
		t.Fatalf("count waiting: %v", err)
	}
	if waitingCount != 1 {
		t.Fatalf("expected 1 mode-queue waiting row, got %d", waitingCount)
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

	if _, err := st.JoinModeQueue(ctx, queueID, userA.ID, ""); err != nil {
		t.Fatalf("join A: %v", err)
	}
	if _, err := st.JoinModeQueue(ctx, queueID, userB.ID, ""); err != nil {
		t.Fatalf("join B: %v", err)
	}

	_, err = st.JoinModeQueue(ctx, queueID, userA.ID, "")
	if err == nil {
		t.Fatal("expected error when joining queue while already matched")
	}
}

func TestJoinModeQueueCompositionMatchmaking(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	slug := "comp-" + uuid.NewString()
	manifest := &gameclient.Manifest{
		Modes: []gameclient.ModeManifest{{
			Key:          "roles",
			DisplayName:  "Roles",
			SeatTemplate: json.RawMessage(`{"Team":{"count":1,"DPS":{"count":1},"Tank":{"count":1}}}`),
		}},
		Status:     gameclient.StatusResponse{Game: "Comp Test", Version: "1.0.0"},
		ETag:       `"comp"`,
		RawJSON:    []byte(`{"modes":[{"key":"roles"}]}`),
		SHA256Hash: uuid.NewString(),
	}
	result, err := st.RegisterGame(ctx, RegisterGameParams{
		Slug:       slug,
		PlayURL:    "https://play.example.com/" + slug,
		APIBaseURL: "https://api.example.com/" + slug,
	}, manifest)
	if err != nil {
		t.Fatalf("RegisterGame: %v", err)
	}
	cleaner.TrackGame(result.Game.ID)

	modes, err := st.ListGameModesByGameID(ctx, result.Game.ID)
	if err != nil {
		t.Fatalf("ListGameModesByGameID: %v", err)
	}
	queues, err := st.ListModeQueuesByModeID(ctx, modes[0].ID)
	if err != nil {
		t.Fatalf("ListModeQueuesByModeID: %v", err)
	}
	queueID := queues[0].ID

	userDPS, err := st.CreateUser(ctx, CreateUserParams{Email: "dps-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser dps: %v", err)
	}
	cleaner.TrackUser(userDPS.ID)
	userTank, err := st.CreateUser(ctx, CreateUserParams{Email: "tank-" + uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("CreateUser tank: %v", err)
	}
	cleaner.TrackUser(userTank.ID)

	first, err := st.JoinModeQueue(ctx, queueID, userDPS.ID, "DPS")
	if err != nil {
		t.Fatalf("join dps: %v", err)
	}
	if first.Status != QueueStatusWaiting {
		t.Fatalf("expected waiting, got %s", first.Status)
	}

	second, err := st.JoinModeQueue(ctx, queueID, userTank.ID, "Tank")
	if err != nil {
		t.Fatalf("join tank: %v", err)
	}
	if second.Status != QueueStatusMatched || second.SessionID == nil {
		t.Fatalf("expected match, got %+v", second)
	}

	participants, err := st.ListSessionSeatAssignments(ctx, *second.SessionID)
	if err != nil {
		t.Fatalf("ListSessionParticipants: %v", err)
	}
	if len(participants) != 2 {
		t.Fatalf("expected 2 participants, got %d", len(participants))
	}
	seatKeys := map[string]struct{}{}
	for _, p := range participants {
		seatKeys[p.SeatKey] = struct{}{}
	}
	if _, ok := seatKeys["Team-1-DPS-1"]; !ok {
		t.Fatalf("expected DPS seat, got %+v", seatKeys)
	}
	if _, ok := seatKeys["Team-1-Tank-1"]; !ok {
		t.Fatalf("expected Tank seat, got %+v", seatKeys)
	}
}
