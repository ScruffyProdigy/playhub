package store

import (
	"testing"

	"github.com/google/uuid"
)

func TestSetUserStarterAvatar(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := t.Context()

	user, err := st.CreateUser(ctx, CreateUserParams{
		Email: "avatar-" + uuid.NewString() + "@example.com",
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(user.ID)

	updated, err := st.SetUserStarterAvatar(ctx, user.ID, "storm", "https://joinquest.cc")
	if err != nil {
		t.Fatalf("SetUserStarterAvatar: %v", err)
	}
	if updated.AvatarKey == nil || *updated.AvatarKey != "storm" {
		t.Fatalf("avatar key: %+v", updated.AvatarKey)
	}
	if updated.AvatarURL == nil || *updated.AvatarURL != "https://joinquest.cc/avatars/storm.png" {
		t.Fatalf("avatar url: %+v", updated.AvatarURL)
	}

	_, err = st.SetUserStarterAvatar(ctx, user.ID, "not-real", "https://joinquest.cc")
	if err != ErrInvalidAvatarKey {
		t.Fatalf("expected ErrInvalidAvatarKey, got %v", err)
	}
}
