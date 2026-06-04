package store

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestConsumeMagicLinkByTokenHashIsSingleUse(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	email := "consume-race-" + uuid.NewString() + "@example.com"
	cleaner.TrackEmail(email)

	tokenHash := "race-hash-" + uuid.NewString()
	link, err := st.CreateMagicLink(ctx, CreateMagicLinkParams{
		Email:     email,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		t.Fatalf("CreateMagicLink: %v", err)
	}

	now := time.Now()
	var wg sync.WaitGroup
	successes := make(chan *MagicLink, 2)
	errors := make(chan error, 2)

	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			consumed, err := st.ConsumeMagicLinkByTokenHash(ctx, tokenHash, now)
			if err != nil {
				errors <- err
				return
			}
			successes <- consumed
		}()
	}
	wg.Wait()
	close(successes)
	close(errors)

	var okCount int
	for consumed := range successes {
		okCount++
		if consumed.ID != link.ID {
			t.Fatalf("unexpected link id %s", consumed.ID)
		}
	}
	if okCount != 1 {
		t.Fatalf("expected exactly one successful consume, got %d", okCount)
	}

	var notFound int
	for err := range errors {
		if err == ErrNotFound {
			notFound++
		}
	}
	if notFound != 1 {
		t.Fatalf("expected one ErrNotFound, got %d errors from channel", notFound)
	}
}
