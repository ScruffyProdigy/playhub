package graph

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/graph/model"
	"github.com/scruffyprodigy/playhub/internal/store"
)

func (r *mutationResolver) joinQueueInternal(ctx context.Context, modeQueueID uuid.UUID, queuePath string, partyInput *model.PartyNodeInput) (*model.JoinResult, error) {
	st, err := r.requireStore()
	if err != nil {
		return nil, err
	}

	userID, err := requireAuthUserID(ctx)
	if err != nil {
		return nil, err
	}

	party, err := toStorePartyInput(ctx, userID, partyInput)
	if err != nil {
		return nil, err
	}

	result, err := st.JoinModeQueue(ctx, modeQueueID, userID, queuePath, party)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyMatched) {
			return nil, fmt.Errorf("you already have an active match; finish or leave your current queue first")
		}
		if errors.Is(err, store.ErrActiveGame) {
			return nil, fmt.Errorf("you have a game in progress — use Leave game at the top of the page")
		}
		if errors.Is(err, store.ErrNotFound) {
			return nil, fmt.Errorf("queue not found")
		}
		return nil, err
	}

	var launchURLs map[uuid.UUID]string
	if result.Status == store.QueueStatusMatched && result.SessionID != nil {
		game, gErr := st.GetGameByID(ctx, result.GameID)
		if gErr != nil {
			return nil, gErr
		}
		launchURLs, err = r.finalizeMatchedSession(ctx, game, *result.SessionID, result.NotifyUserIDs)
		if err != nil {
			return nil, err
		}
	}

	if result.SwitchedFrom != nil {
		if err := r.publishQueueSwitchFrom(ctx, st, result.SwitchedFrom, userID); err != nil {
			return nil, err
		}
	}

	if err := r.publishQueueResult(ctx, result, launchURLs); err != nil {
		return nil, err
	}

	return joinResultAfterJoin(result, userID, launchURLs), nil
}
