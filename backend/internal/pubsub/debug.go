package pubsub

import (
	"log"
	"os"
	"strings"
)

// DebugEnabled reports whether queue pub/sub tracing is on (LOBBY_PUBSUB_DEBUG=1).
func DebugEnabled() bool {
	v := strings.TrimSpace(os.Getenv("LOBBY_PUBSUB_DEBUG"))
	return v == "1" || strings.EqualFold(v, "true")
}

// DebugLog writes a pubsub trace line when LOBBY_PUBSUB_DEBUG is enabled.
func DebugLog(format string, args ...any) {
	if DebugEnabled() {
		log.Printf("pubsub: "+format, args...)
	}
}
