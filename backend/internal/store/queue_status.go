package store

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// UserQueueView is the current queue/match state for a user and game.
type UserQueueView struct {
	InQueue     bool
	Waiting     bool
	Matched     bool
	QueuedCount int
	SessionID   *uuid.UUID
}

// GetUserQueueView returns whether the user is waiting, matched, or not in queue.
func (s *Store) GetUserQueueView(ctx context.Context, gameID, userID uuid.UUID) (*UserQueueView, error) {
	if err := s.ExpireStaleMatchedQueue(ctx, gameID, userID); err != nil {
		return nil, err
	}

	if _, err := s.GetWaitingQueueEntry(ctx, gameID, userID); err == nil {
		count, err := s.countWaitingQueueEntries(ctx, gameID)
		if err != nil {
			return nil, err
		}
		return &UserQueueView{
			InQueue:     true,
			Waiting:     true,
			QueuedCount: count,
		}, nil
	} else if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	session, err := s.GetMatchedSessionForUserAndGame(ctx, gameID, userID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	if session != nil {
		return &UserQueueView{
			InQueue:   true,
			Matched:   true,
			SessionID: &session.ID,
		}, nil
	}

	return &UserQueueView{}, nil
}
