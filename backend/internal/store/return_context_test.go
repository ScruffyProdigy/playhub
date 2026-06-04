package store

import (
	"testing"

	"github.com/google/uuid"
)

func TestReturnContextValidate(t *testing.T) {
	gameID := uuid.New()
	queueID := uuid.New()
	ctx := CatalogLFGReturnContext(gameID, queueID)
	if err := ctx.Validate(); err != nil {
		t.Fatalf("Validate() error: %v", err)
	}
	if ctx.Kind != ReturnKindCatalogLFG || ctx.Path != "/" {
		t.Fatalf("unexpected context: %+v", ctx)
	}
}

func TestReturnContextRejectsOpenRedirect(t *testing.T) {
	ctx := ReturnContext{Kind: ReturnKindCatalogLFG, Path: "//evil.example"}
	if err := ctx.Validate(); err == nil {
		t.Fatal("expected validation error for open redirect path")
	}
}

func TestDecodeReturnContextDefaults(t *testing.T) {
	ctx, err := decodeReturnContext(nil)
	if err != nil {
		t.Fatalf("decodeReturnContext() error: %v", err)
	}
	if ctx.Path != "/" || ctx.Kind != ReturnKindCatalogLFG {
		t.Fatalf("defaults: %+v", ctx)
	}
}
