package spiritanimal

import (
	"strings"
)

// UserFacingError maps internal failures to safe client-facing text.
func UserFacingError(err error) string {
	if err == nil {
		return "Something went wrong while summoning your mascots."
	}
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "quota"), strings.Contains(msg, "429"), strings.Contains(msg, "billing"):
		return "OpenAI quota exceeded — add billing or credits, then try again."
	case strings.Contains(msg, "permission denied"):
		return "Could not save mascot images on the server. Try again in a moment."
	case strings.Contains(msg, "expected 5"), strings.Contains(msg, "invalid questions"), strings.Contains(msg, "invalid totems"):
		return "We could not parse the reading. Tap Start over to draw a fresh one."
	case strings.Contains(msg, "invalid answer"):
		return "One of those answers was not valid for this reading. Try again."
	case strings.Contains(msg, "timed out"), strings.Contains(msg, "deadline exceeded"):
		return "Mascot generation is taking longer than expected. Try again in a moment."
	case strings.Contains(msg, "http 500"), strings.Contains(msg, "http 502"), strings.Contains(msg, "http 503"),
		strings.Contains(msg, "server_error"), strings.Contains(msg, "upstream connect error"):
		return "OpenAI had a temporary hiccup while drawing your mascots. Tap Start over to try again."
	case strings.Contains(msg, "dall-e-3"), strings.Contains(msg, "dall-e-2"):
		return "That reading used a retired image model. Tap Start over to start fresh."
	default:
		trimmed := strings.TrimSpace(err.Error())
		if trimmed == "" {
			return "Something went wrong while summoning your mascots."
		}
		return trimmed
	}
}
