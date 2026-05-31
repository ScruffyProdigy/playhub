package store

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestCleanerRemovesCreatedFixtures(t *testing.T) {
	st := openTestStore(t)
	ctx := context.Background()

	game, err := st.InsertTestGame(ctx, "Cleanup Test Game "+uuid.NewString())
	if err != nil {
		t.Fatalf("InsertTestGame failed: %v", err)
	}

	user, err := st.CreateUser(ctx, CreateUserParams{
		Email: "cleanup-" + uuid.NewString() + "@example.com",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	cleaner := &TestCleaner{st: st}
	cleaner.TrackGame(game.ID)
	cleaner.TrackUser(user.ID)
	cleaner.run(ctx)

	if _, err := st.GetGameByID(ctx, game.ID); err != ErrNotFound {
		t.Fatalf("expected game %s to be deleted, got err=%v", game.ID, err)
	}
	if _, err := st.GetUserByID(ctx, user.ID); err != ErrNotFound {
		t.Fatalf("expected user %s to be deleted, got err=%v", user.ID, err)
	}
}
