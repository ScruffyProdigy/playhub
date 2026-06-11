package formingworker

import (
	"testing"
	"time"
)

func TestProvisionBackoff(t *testing.T) {
	if ProvisionBackoff(0) != 100*time.Millisecond {
		t.Fatalf("attempt 0")
	}
	if ProvisionBackoff(1) != 500*time.Millisecond {
		t.Fatalf("attempt 1")
	}
	if ProvisionBackoff(5) != 30*time.Second {
		t.Fatalf("attempt 5 cap")
	}
}
