package graph

import (
	"context"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/graph/model"
	"github.com/scruffyprodigy/playhub/internal/store"
)

func loadAuthorizedGameSession(ctx context.Context, st *store.Store, matchID string) (uuid.UUID, *store.Session, error) {
	sessionID, err := parseUUID(matchID, "match id")
	if err != nil {
		return uuid.Nil, nil, err
	}
	session, err := st.GetSessionByID(ctx, sessionID)
	if err != nil {
		return uuid.Nil, nil, err
	}
	if err := requireGameServiceForSessionGame(ctx, session.GameID); err != nil {
		return uuid.Nil, nil, err
	}
	return sessionID, session, nil
}

func returnDestinationFromContext(ctx store.ReturnContext) *model.ReturnDestination {
	defaultDest := &model.ReturnDestination{Path: "/", Kind: store.ReturnKindCatalogLFG}
	if err := ctx.Validate(); err != nil {
		return defaultDest
	}
	return &model.ReturnDestination{Path: ctx.Path, Kind: ctx.Kind}
}
