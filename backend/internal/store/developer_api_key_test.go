package store

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestDeveloperAPIKeyRoundTrip(t *testing.T) {
	os.Setenv("MAGIC_LINK_PEPPER", "test-pepper")
	t.Cleanup(func() { os.Unsetenv("MAGIC_LINK_PEPPER") })

	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	user, err := st.CreateUser(ctx, CreateUserParams{
		Email:       "dev-key-" + uuid.NewString() + "@example.com",
		DisplayName: "Dev Key User",
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(user.ID)

	raw, key, err := st.CreateDeveloperAPIKey(ctx, user.ID, "CI")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !strings.HasPrefix(raw, "lq_dev_") {
		t.Fatalf("expected developer api key prefix, got %q", raw)
	}

	got, err := st.VerifyDeveloperAPIKey(ctx, raw)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if got != user.ID {
		t.Fatalf("user id mismatch: got %v want %v", got, user.ID)
	}

	keys, err := st.ListDeveloperAPIKeys(ctx, user.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(keys) != 1 || keys[0].ID != key.ID {
		t.Fatalf("unexpected keys: %+v", keys)
	}

	if err := st.RevokeDeveloperAPIKey(ctx, key.ID, user.ID); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if _, err := st.VerifyDeveloperAPIKey(ctx, raw); err == nil {
		t.Fatal("expected revoked key to fail verification")
	}
}

func TestRevokeDeveloperAPIKeyWrongOwner(t *testing.T) {
	os.Setenv("MAGIC_LINK_PEPPER", "test-pepper")
	t.Cleanup(func() { os.Unsetenv("MAGIC_LINK_PEPPER") })

	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	owner, err := st.CreateUser(ctx, CreateUserParams{
		Email:       "owner-" + uuid.NewString() + "@example.com",
		DisplayName: "Owner",
	})
	if err != nil {
		t.Fatalf("CreateUser owner: %v", err)
	}
	cleaner.TrackUser(owner.ID)

	other, err := st.CreateUser(ctx, CreateUserParams{
		Email:       "other-" + uuid.NewString() + "@example.com",
		DisplayName: "Other",
	})
	if err != nil {
		t.Fatalf("CreateUser other: %v", err)
	}
	cleaner.TrackUser(other.ID)

	_, key, err := st.CreateDeveloperAPIKey(ctx, owner.ID, "Mine")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := st.RevokeDeveloperAPIKey(ctx, key.ID, other.ID); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}
