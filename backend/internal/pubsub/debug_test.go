package pubsub

import (
	"os"
	"testing"
)

func TestDebugEnabled(t *testing.T) {
	t.Setenv("LOBBY_PUBSUB_DEBUG", "")
	if DebugEnabled() {
		t.Fatal("expected disabled when unset")
	}
	t.Setenv("LOBBY_PUBSUB_DEBUG", "true")
	if !DebugEnabled() {
		t.Fatal("expected enabled for true")
	}
	t.Setenv("LOBBY_PUBSUB_DEBUG", "0")
	if DebugEnabled() {
		t.Fatal("expected disabled for 0")
	}
}

func TestDebugLogNoPanicWhenDisabled(t *testing.T) {
	os.Unsetenv("LOBBY_PUBSUB_DEBUG")
	DebugLog("test %s", "message")
}
