package graph

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/gameclient"
	"github.com/scruffyprodigy/playhub/internal/pubsub"
	"github.com/scruffyprodigy/playhub/internal/store"
)

var sessionProvisionLocks sync.Map // session ID string -> *sync.Mutex

func lockSessionProvision(sessionID uuid.UUID) func() {
	v, _ := sessionProvisionLocks.LoadOrStore(sessionID.String(), &sync.Mutex{})
	mu := v.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

// HandleFormingReconciled provisions a fired match and publishes queue/table events.
func (r *Resolver) HandleFormingReconciled(ctx context.Context, result *store.FormingReconcileResult) error {
	if r == nil || result == nil || !result.Fired || result.SessionID == nil {
		return nil
	}
	return r.provisionMatchedSession(ctx, matchedSessionWork{
		GameID:        result.GameID,
		ModeQueueID:   result.ModeQueueID,
		SessionID:     *result.SessionID,
		NotifyUserIDs: result.NotifyUserIDs,
		TableIDs:      result.TableIDs,
	})
}

// HandleUnprovisionedSession retries game provision for an already-matched session.
func (r *Resolver) HandleUnprovisionedSession(ctx context.Context, pending store.UnprovisionedSession) error {
	if r == nil || pending.SessionID == uuid.Nil {
		return nil
	}
	return r.provisionMatchedSession(ctx, matchedSessionWork{
		GameID:        pending.GameID,
		ModeQueueID:   pending.ModeQueueID,
		SessionID:     pending.SessionID,
		NotifyUserIDs: pending.NotifyUserIDs,
	})
}

type matchedSessionWork struct {
	GameID        uuid.UUID
	ModeQueueID   uuid.UUID
	SessionID     uuid.UUID
	NotifyUserIDs []uuid.UUID
	TableIDs      []uuid.UUID
}

func (r *Resolver) provisionMatchedSession(ctx context.Context, work matchedSessionWork) error {
	st, err := r.requireStore()
	if err != nil {
		return err
	}

	game, err := st.GetGameByID(ctx, work.GameID)
	if err != nil {
		return err
	}

	launchURLs, err := r.finalizeMatchedSession(ctx, game, work.SessionID, work.NotifyUserIDs)
	if err != nil {
		var banned *gameclient.BannedPlayersError
		if errors.As(err, &banned) {
			if rbErr := st.RollbackMatchedSession(ctx, work.SessionID, parseBannedLobbyUserIDs(banned.BannedLobbyUserIDs)); rbErr != nil {
				return err
			}
			return err
		}
		log.Printf("handoff: provision deferred session=%s queue=%s: %v", work.SessionID, work.ModeQueueID, err)
		if r.FormingWorker != nil {
			r.FormingWorker.ScheduleProvisionRetry(store.UnprovisionedSession{
				SessionID:     work.SessionID,
				GameID:        work.GameID,
				ModeQueueID:   work.ModeQueueID,
				NotifyUserIDs: work.NotifyUserIDs,
			})
		}
		return nil
	}

	if r.FormingWorker != nil {
		r.FormingWorker.ClearProvisionRetry(work.SessionID)
	}

	for _, tableID := range work.TableIDs {
		table, tableErr := st.GetRoomTableByID(ctx, tableID)
		if tableErr != nil {
			return tableErr
		}
		if err := r.publishTableSeatStarted(ctx, tableID, work.NotifyUserIDs, launchURLs); err != nil {
			return err
		}
		if err := r.publishTableUpdated(ctx, table.RoomID, tableID); err != nil {
			return err
		}
	}

	queueResult := &store.QueueJoinResult{
		GameID:        work.GameID,
		ModeQueueID:   work.ModeQueueID,
		Status:        store.QueueStatusMatched,
		SessionID:     &work.SessionID,
		NotifyUserIDs: work.NotifyUserIDs,
	}
	pubsub.DebugLog(
		"forming matched queue=%s session=%s notifyUsers=%d",
		work.ModeQueueID,
		work.SessionID,
		len(work.NotifyUserIDs),
	)
	return r.publishQueueResult(ctx, queueResult, launchURLs)
}

func (r *Resolver) scheduleFormingReconcile(modeQueueID uuid.UUID) {
	if r == nil || r.FormingWorker == nil || modeQueueID == uuid.Nil {
		return
	}
	r.FormingWorker.Schedule(modeQueueID)
}
