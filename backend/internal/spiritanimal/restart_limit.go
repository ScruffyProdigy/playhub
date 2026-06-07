package spiritanimal

import "time"

const (
	// MaxReadingStartsPerHour limits new journeys and force-restarts to control OpenAI spend.
	MaxReadingStartsPerHour = 4
)

// ReadingStartWindow is how far back start counts are measured.
func ReadingStartWindow() time.Duration {
	return time.Hour
}
