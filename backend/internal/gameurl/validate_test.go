package gameurl

import (
	"context"
	"strings"
	"testing"
)

func TestValidateOutboundURLLocalhostDev(t *testing.T) {
	ctx := context.Background()
	if err := ValidateOutboundURL(ctx, "http://localhost:3001", false); err != nil {
		t.Fatalf("localhost http in dev: %v", err)
	}
}

func TestValidateOutboundURLRequiresHTTPSInProduction(t *testing.T) {
	ctx := context.Background()
	err := ValidateOutboundURL(ctx, "http://example.com", true)
	if err == nil || !strings.Contains(err.Error(), "HTTPS") {
		t.Fatalf("expected HTTPS error, got %v", err)
	}
}

func TestValidateOutboundURLBlocksPrivateIP(t *testing.T) {
	ctx := context.Background()
	err := ValidateOutboundURL(ctx, "https://127.0.0.1.nip.io", true)
	if err == nil {
		t.Fatal("expected blocked address error for loopback host")
	}
}
