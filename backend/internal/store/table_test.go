package store

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/gameclient"
)

func TestTableSitStart(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	host, err := st.CreateUser(ctx, CreateUserParams{Email: "table-host@example.com"})
	if err != nil {
		t.Fatalf("create host: %v", err)
	}
	cleaner.TrackUser(host.ID)
	guest, err := st.CreateUser(ctx, CreateUserParams{Email: "table-guest@example.com"})
	if err != nil {
		t.Fatalf("create guest: %v", err)
	}
	cleaner.TrackUser(guest.ID)

	game, mode, modeQueueID := setupWordHuntMode(t, st, cleaner)
	_ = modeQueueID

	room, err := st.CreateRoom(ctx, host.ID)
	if err != nil {
		t.Fatalf("create room: %v", err)
	}
	if err := st.addRoomMemberDirect(ctx, room.ID, guest.ID); err != nil {
		t.Fatalf("add guest: %v", err)
	}

	table, err := st.CreateTable(ctx, room.ID, game.ID, mode.ID, host.ID)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}

	seats, err := st.ListGameModeSeats(ctx, mode.ID)
	if err != nil {
		t.Fatalf("list seats: %v", err)
	}
	if len(seats) < 4 {
		t.Fatalf("need at least 4 seats, got %d", len(seats))
	}

	clueKeys := make([]string, 0, 2)
	guessKeys := make([]string, 0, 2)
	for _, seat := range seats {
		path := seatQueuePathValue(seat)
		switch path {
		case "ClueGiver":
			if len(clueKeys) < 2 {
				clueKeys = append(clueKeys, seat.SeatKey)
			}
		case "Guesser":
			if len(guessKeys) < 2 {
				guessKeys = append(guessKeys, seat.SeatKey)
			}
		}
	}
	if len(clueKeys) < 2 || len(guessKeys) < 2 {
		t.Fatalf("missing clue/guesser seats: clue=%v guess=%v", clueKeys, guessKeys)
	}

	if _, err := st.SitAtTable(ctx, table.ID, host.ID, clueKeys[0]); err != nil {
		t.Fatalf("host sit: %v", err)
	}
	if _, err := st.SitAtTable(ctx, table.ID, guest.ID, clueKeys[1]); err != nil {
		t.Fatalf("guest sit clue: %v", err)
	}

	third, err := st.CreateUser(ctx, CreateUserParams{Email: "table-third@example.com"})
	if err != nil {
		t.Fatalf("create third: %v", err)
	}
	cleaner.TrackUser(third.ID)
	if err := st.addRoomMemberDirect(ctx, room.ID, third.ID); err != nil {
		t.Fatalf("add third: %v", err)
	}
	if _, err := st.SitAtTable(ctx, table.ID, third.ID, guessKeys[0]); err != nil {
		t.Fatalf("third sit: %v", err)
	}

	canStart, err := st.TableCanStart(ctx, table.ID)
	if err != nil {
		t.Fatalf("canStart: %v", err)
	}
	if canStart {
		t.Fatal("canStart should be false with only 3 players")
	}

	fourth, err := st.CreateUser(ctx, CreateUserParams{Email: "table-fourth@example.com"})
	if err != nil {
		t.Fatalf("create fourth: %v", err)
	}
	cleaner.TrackUser(fourth.ID)
	if err := st.addRoomMemberDirect(ctx, room.ID, fourth.ID); err != nil {
		t.Fatalf("add fourth: %v", err)
	}
	if _, err := st.SitAtTable(ctx, table.ID, fourth.ID, guessKeys[1]); err != nil {
		t.Fatalf("fourth sit: %v", err)
	}

	canStart, err = st.TableCanStart(ctx, table.ID)
	if err != nil {
		t.Fatalf("canStart: %v", err)
	}
	if !canStart {
		t.Fatal("canStart should be true with full roster")
	}

	king, err := st.TableKingUserID(ctx, table.ID)
	if err != nil || king == nil || *king != host.ID {
		t.Fatalf("king = %v, want host %s", king, host.ID)
	}

	result, err := st.StartTable(ctx, table.ID, host.ID)
	if err != nil {
		t.Fatalf("start table: %v", err)
	}
	if result.SessionID == uuid.Nil {
		t.Fatal("expected session id")
	}

	participants, err := st.ListSessionSeatAssignments(ctx, result.SessionID)
	if err != nil {
		t.Fatalf("participants: %v", err)
	}
	if len(participants) != 4 {
		t.Fatalf("participant count = %d, want 4", len(participants))
	}
}

func TestTableQueueExclusion(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	user, err := st.CreateUser(ctx, CreateUserParams{Email: "table-queue@example.com"})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	cleaner.TrackUser(user.ID)

	game, mode, modeQueueID := setupWordHuntMode(t, st, cleaner)

	table, err := st.CreatePrivateTable(ctx, user.ID, game.ID, mode.ID)
	if err != nil {
		t.Fatalf("create private table: %v", err)
	}
	seats, err := st.ListGameModeSeats(ctx, mode.ID)
	if err != nil {
		t.Fatalf("list seats: %v", err)
	}
	if _, err := st.SitAtTable(ctx, table.ID, user.ID, seats[0].SeatKey); err != nil {
		t.Fatalf("sit: %v", err)
	}

	if _, err := st.JoinModeQueue(ctx, modeQueueID, user.ID, "ClueGiver"); err != nil {
		t.Fatalf("join queue: %v", err)
	}

	seat, err := st.GetUserTableSeat(ctx, user.ID)
	if err != nil {
		t.Fatalf("get table seat: %v", err)
	}
	if seat != nil {
		t.Fatal("expected table seat cleared after join queue")
	}
}

func setupWordHuntMode(t *testing.T, st *Store, cleaner *TestCleaner) (*Game, *GameMode, uuid.UUID) {
	t.Helper()
	ctx := context.Background()
	slug := "table-wh-" + uuid.NewString()
	manifest := &gameclient.Manifest{
		Modes: []gameclient.ModeManifest{{
			Key:         "party",
			DisplayName: "Word Hunt Party",
			Min:         4,
			Max:         9,
			SeatTemplate: json.RawMessage(`{
				"ClueGiver":{"displayName":"Clue Giver","name":["Red","Blue","Green"],"min":2,"max":3,"sizeForQueue":2},
				"Guesser":{"count":6,"min":2,"max":6,"sizeForQueue":4}
			}`),
		}},
		Status:     gameclient.StatusResponse{Game: "Word Hunt", Version: "1.0.0"},
		ETag:       `"wh"`,
		RawJSON:    []byte(`{"modes":[{"key":"party"}]}`),
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
	return result.Game, &modes[0], queues[0].ID
}
func (s *Store) addRoomMemberDirect(ctx context.Context, roomID, userID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO room_members (room_id, user_id) VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, roomID, userID)
	return err
}
