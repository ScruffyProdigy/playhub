package spiritanimal

import (
	"testing"
	"time"
)

func TestJourneyEligibilityAtAdminBypass(t *testing.T) {
	last := time.Now().Add(-24 * time.Hour)
	got := JourneyEligibilityAt(&last, time.Now(), true)
	if !got.CanBegin {
		t.Fatal("expected admin bypass")
	}
}

func TestJourneyEligibilityAtCooldownActive(t *testing.T) {
	last := time.Now().Add(-5 * 24 * time.Hour)
	got := JourneyEligibilityAt(&last, time.Now(), false)
	if got.CanBegin {
		t.Fatal("expected cooldown")
	}
	if got.DaysRemaining <= 0 || got.CooldownEndsAt == nil {
		t.Fatalf("expected remaining days, got %+v", got)
	}
}

func TestJourneyEligibilityAtExpired(t *testing.T) {
	last := time.Now().Add(-31 * 24 * time.Hour)
	got := JourneyEligibilityAt(&last, time.Now(), false)
	if !got.CanBegin {
		t.Fatal("expected cooldown expired")
	}
}
