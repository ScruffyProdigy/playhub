package spiritanimal

import (
	"math"
	"time"
)

const JourneyCooldownDays = 30

// JourneyCooldown is the wait between completed spirit animal journeys.
func JourneyCooldown() time.Duration {
	return JourneyCooldownDays * 24 * time.Hour
}

// JourneyEligibility describes whether a user may start a new journey.
type JourneyEligibility struct {
	CanBegin       bool
	CooldownEndsAt *time.Time
	DaysRemaining  int
}

// JourneyEligibilityAt reports journey availability from the last completion time.
func JourneyEligibilityAt(lastCompleted *time.Time, now time.Time, isAdmin bool) JourneyEligibility {
	if isAdmin {
		return JourneyEligibility{CanBegin: true}
	}
	if lastCompleted == nil {
		return JourneyEligibility{CanBegin: true}
	}

	endsAt := lastCompleted.Add(JourneyCooldown())
	if !now.Before(endsAt) {
		return JourneyEligibility{CanBegin: true}
	}

	remaining := endsAt.Sub(now)
	days := int(math.Ceil(remaining.Hours() / 24))
	if days < 1 {
		days = 1
	}
	return JourneyEligibility{
		CanBegin:       false,
		CooldownEndsAt: &endsAt,
		DaysRemaining:  days,
	}
}
