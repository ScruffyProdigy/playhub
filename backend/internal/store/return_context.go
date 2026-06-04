package store

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Return context kinds (extend with "room", "party", etc.).
const (
	ReturnKindCatalogLFG = "catalog_lfg"
)

// ReturnContext tells the Lobby return hub where to send a player after a match.
// Path must be a same-origin relative path (e.g. "/").
type ReturnContext struct {
	Kind        string `json:"kind"`
	Path        string `json:"path"`
	GameID      string `json:"gameId,omitempty"`
	ModeQueueID string `json:"modeQueueId,omitempty"`
}

// CatalogLFGReturnContext is the default when joining a catalog mode queue.
func CatalogLFGReturnContext(gameID, modeQueueID uuid.UUID) ReturnContext {
	return ReturnContext{
		Kind:        ReturnKindCatalogLFG,
		Path:        "/",
		GameID:      gameID.String(),
		ModeQueueID: modeQueueID.String(),
	}
}

func (c ReturnContext) Validate() error {
	if strings.TrimSpace(c.Kind) == "" {
		return fmt.Errorf("store: return context kind is required")
	}
	path := strings.TrimSpace(c.Path)
	if path == "" {
		return fmt.Errorf("store: return context path is required")
	}
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") {
		return fmt.Errorf("store: return context path must be a relative path starting with /")
	}
	return nil
}

func encodeReturnContext(ctx ReturnContext) ([]byte, error) {
	if err := ctx.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(ctx)
}

func decodeReturnContext(raw []byte) (ReturnContext, error) {
	if len(raw) == 0 {
		return ReturnContext{Kind: ReturnKindCatalogLFG, Path: "/"}, nil
	}
	var ctx ReturnContext
	if err := json.Unmarshal(raw, &ctx); err != nil {
		return ReturnContext{}, fmt.Errorf("store: decode return context: %w", err)
	}
	if strings.TrimSpace(ctx.Path) == "" {
		ctx.Path = "/"
	}
	if strings.TrimSpace(ctx.Kind) == "" {
		ctx.Kind = ReturnKindCatalogLFG
	}
	return ctx, nil
}
