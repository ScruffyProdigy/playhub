package graph

import (
	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/graph/model"
	"github.com/scruffyprodigy/playhub/internal/store"
)

// joinResultFromQueueJoin is the joinGame mutation response: always "in queue" from the
// client's perspective. Match details are delivered via queueUpdated / myQueueStatus.
func joinResultFromQueueJoin(result *store.QueueJoinResult) *model.JoinResult {
	if result == nil {
		return &model.JoinResult{Queued: false}
	}

	count := result.QueuedCount
	return &model.JoinResult{
		Queued:      true,
		QueuedCount: &count,
	}
}

// joinResultAfterJoin is the joinGame response. Waiting users get queued only; the user
// who completes a match also receives sessionId/joinUrl (same as myQueueStatus).
func joinResultAfterJoin(result *store.QueueJoinResult, userID uuid.UUID, launchURLs map[uuid.UUID]string) *model.JoinResult {
	if result == nil {
		return &model.JoinResult{Queued: false}
	}

	if result.Status == store.QueueStatusMatched && result.SessionID != nil {
		sessionID := result.SessionID.String()
		resp := &model.JoinResult{
			Queued:    false,
			SessionID: &sessionID,
		}
		if launchURLs != nil {
			if url, ok := launchURLs[userID]; ok && url != "" {
				resp.JoinURL = &url
			}
		}
		return resp
	}

	return joinResultFromQueueJoin(result)
}

func joinResultFromQueueView(view *store.UserQueueView, launchURL string) *model.JoinResult {
	if view == nil || !view.InQueue {
		return &model.JoinResult{Queued: false}
	}

	if view.Waiting {
		count := view.QueuedCount
		return &model.JoinResult{
			Queued:      true,
			QueuedCount: &count,
		}
	}

	if view.Matched && view.SessionID != nil {
		sessionID := view.SessionID.String()
		result := &model.JoinResult{
			Queued:    false,
			SessionID: &sessionID,
		}
		if launchURL != "" {
			result.JoinURL = &launchURL
		}
		return result
	}

	return &model.JoinResult{Queued: false}
}
