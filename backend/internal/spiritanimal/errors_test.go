package spiritanimal

import (
	"errors"
	"strings"
	"testing"
)

func TestUserFacingErrorMapsOpenAIOutages(t *testing.T) {
	got := UserFacingError(errors.New("openai image: HTTP 503: upstream connect error"))
	if !strings.Contains(got, "temporary hiccup") {
		t.Fatalf("expected friendly outage message, got %q", got)
	}
}

func TestUserFacingErrorPreservesCooldown(t *testing.T) {
	msg := "Your next spirit animal journey opens in 12 days"
	got := UserFacingError(errors.New(msg))
	if got != msg {
		t.Fatalf("expected %q, got %q", msg, got)
	}
}
