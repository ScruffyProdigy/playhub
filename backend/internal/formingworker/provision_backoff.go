package formingworker

import "time"

// ProvisionBackoff returns the delay before a provision retry attempt (0-based).
func ProvisionBackoff(attempt int) time.Duration {
	switch attempt {
	case 0:
		return 100 * time.Millisecond
	case 1:
		return 500 * time.Millisecond
	case 2:
		return 2 * time.Second
	case 3:
		return 5 * time.Second
	case 4:
		return 15 * time.Second
	default:
		return 30 * time.Second
	}
}
