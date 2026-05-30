package graph

import (
	"context"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/graph/model"
	"github.com/scruffyprodigy/playhub/internal/pubsub"
	"github.com/scruffyprodigy/playhub/internal/store"
)

func (r *Resolver) publishQueueResult(ctx context.Context, gameID uuid.UUID, result *store.QueueJoinResult, launchURLs map[uuid.UUID]string) error {
	if r.PubSub == nil || result == nil {
		return nil
	}

	seen := make(map[uuid.UUID]struct{}, len(result.NotifyUserIDs))
	for _, userID := range result.NotifyUserIDs {
		if _, ok := seen[userID]; ok {
			continue
		}
		seen[userID] = struct{}{}

		event := pubsub.QueueEvent{
			GameID:      gameID.String(),
			QueuedCount: result.QueuedCount,
		}

		switch result.Status {
		case store.QueueStatusMatched:
			event.Status = pubsub.QueueStatusMatched
			if result.SessionID != nil {
				event.SessionID = result.SessionID.String()
				if launchURLs != nil {
					event.JoinURL = launchURLs[userID]
				}
			}
		default:
			event.Status = pubsub.QueueStatusWaiting
		}

		if err := pubsub.PublishQueueEvent(ctx, r.PubSub, userID.String(), event); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) publishQueueLeft(ctx context.Context, gameID, userID uuid.UUID, queuedCount int) error {
	if r.PubSub == nil {
		return nil
	}
	return pubsub.PublishQueueEvent(ctx, r.PubSub, userID.String(), pubsub.QueueEvent{
		GameID:      gameID.String(),
		Status:      pubsub.QueueStatusLeft,
		QueuedCount: queuedCount,
	})
}

func queueUpdateFromView(view *store.UserQueueView, gameID uuid.UUID, launchURL string) *model.QueueUpdate {
	if view == nil || !view.InQueue {
		return nil
	}

	gameIDStr := gameID.String()
	if view.Matched && view.SessionID != nil {
		update := &model.QueueUpdate{
			GameID:      gameIDStr,
			Status:      model.QueueStatusMatched,
			QueuedCount: 0,
		}
		sessionID := view.SessionID.String()
		update.SessionID = &sessionID
		if launchURL != "" {
			update.JoinURL = &launchURL
		}
		return update
	}

	if view.Waiting {
		return &model.QueueUpdate{
			GameID:      gameIDStr,
			Status:      model.QueueStatusWaiting,
			QueuedCount: view.QueuedCount,
		}
	}

	return nil
}

func toGraphQLQueueUpdate(event pubsub.QueueEvent) *model.QueueUpdate {
	update := &model.QueueUpdate{
		GameID:      event.GameID,
		QueuedCount: event.QueuedCount,
	}
	switch event.Status {
	case pubsub.QueueStatusMatched:
		update.Status = model.QueueStatusMatched
	case pubsub.QueueStatusLeft:
		update.Status = model.QueueStatusLeft
	default:
		update.Status = model.QueueStatusWaiting
	}
	if event.SessionID != "" {
		update.SessionID = &event.SessionID
	}
	if event.JoinURL != "" {
		update.JoinURL = &event.JoinURL
	}
	return update
}
