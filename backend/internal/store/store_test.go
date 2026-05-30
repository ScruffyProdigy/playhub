package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

func openTestStore(t *testing.T) *Store {
	t.Helper()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping store integration test")
	}

	db, err := openDB(databaseURL)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	return New(db)
}

func TestStoreUserGameFlow(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	email := "store-test-" + uuid.NewString() + "@example.com"
	user, err := st.CreateUser(ctx, CreateUserParams{
		Email:       email,
		DisplayName: "Store Test User",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	cleaner.TrackUser(user.ID)

	found, err := st.GetUserByEmail(ctx, email)
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if found.ID != user.ID {
		t.Fatalf("expected user id %s, got %s", user.ID, found.ID)
	}

	game, err := st.CreateGame(ctx, "Store Test Game")
	if err != nil {
		t.Fatalf("CreateGame failed: %v", err)
	}
	cleaner.TrackGame(game.ID)

	games, err := st.ListGames(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListGames failed: %v", err)
	}
	if len(games) == 0 {
		t.Fatal("expected at least one game")
	}

	entry, err := st.EnqueueGame(ctx, game.ID, user.ID)
	if err != nil {
		t.Fatalf("EnqueueGame failed: %v", err)
	}
	if entry.Status != "waiting" {
		t.Fatalf("expected waiting queue status, got %s", entry.Status)
	}

	if err := st.LeaveQueue(ctx, game.ID, user.ID); err != nil {
		t.Fatalf("LeaveQueue failed: %v", err)
	}
}

func TestStoreMagicLinkFlow(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	email := "magic-" + uuid.NewString() + "@example.com"
	cleaner.TrackEmail(email)

	token := uuid.NewString()
	link, err := st.CreateMagicLink(ctx, CreateMagicLinkParams{
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		t.Fatalf("CreateMagicLink failed: %v", err)
	}
	if !st.IsMagicLinkValid(link, time.Now()) {
		t.Fatal("expected magic link to be valid")
	}

	found, err := st.GetMagicLinkByToken(ctx, token)
	if err != nil {
		t.Fatalf("GetMagicLinkByToken failed: %v", err)
	}
	if found.ID != link.ID {
		t.Fatalf("expected link id %s, got %s", link.ID, found.ID)
	}

	if err := st.MarkMagicLinkUsed(ctx, link.ID); err != nil {
		t.Fatalf("MarkMagicLinkUsed failed: %v", err)
	}

	found, err = st.GetMagicLinkByToken(ctx, token)
	if err != nil {
		t.Fatalf("GetMagicLinkByToken after use failed: %v", err)
	}
	if st.IsMagicLinkValid(found, time.Now()) {
		t.Fatal("expected magic link to be invalid after use")
	}
}

func TestStoreInventoryFlow(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	user, err := st.CreateUser(ctx, CreateUserParams{
		Email:       "inventory-" + uuid.NewString() + "@example.com",
		DisplayName: "Inventory User",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	cleaner.TrackUser(user.ID)

	game, err := st.CreateGame(ctx, "Inventory Game")
	if err != nil {
		t.Fatalf("CreateGame failed: %v", err)
	}
	cleaner.TrackGame(game.ID)

	description := "A test item"
	good, err := st.CreateDigitalGood(ctx, "Test Good", &description, &game.ID)
	if err != nil {
		t.Fatalf("CreateDigitalGood failed: %v", err)
	}

	if err := st.GrantInventoryItem(ctx, user.ID, good.ID, 2); err != nil {
		t.Fatalf("GrantInventoryItem failed: %v", err)
	}

	items, err := st.ListUserInventory(ctx, user.ID, &game.ID)
	if err != nil {
		t.Fatalf("ListUserInventory failed: %v", err)
	}
	if len(items) != 1 || items[0].Quantity != 2 {
		t.Fatalf("unexpected inventory: %+v", items)
	}

	if err := st.RevokeInventoryItem(ctx, user.ID, good.ID, 1); err != nil {
		t.Fatalf("RevokeInventoryItem failed: %v", err)
	}

	items, err = st.ListUserInventory(ctx, user.ID, &game.ID)
	if err != nil {
		t.Fatalf("ListUserInventory failed: %v", err)
	}
	if len(items) != 1 || items[0].Quantity != 1 {
		t.Fatalf("expected quantity 1 after revoke, got %+v", items)
	}
}
