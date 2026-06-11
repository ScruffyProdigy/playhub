package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type formingFireResult struct {
	session   *Session
	notifyIDs []uuid.UUID
	tableIDs  []uuid.UUID
}

func (s *Store) joinModeQueueFormingTx(
	ctx context.Context,
	tx *sql.Tx,
	joinCtx *ModeQueueJoinContext,
	modeQueueID, callerID uuid.UUID,
	queuePath string,
	partyInput *JoinPartyInput,
) (*QueueJoinResult, error) {
	var members []JoinPartyMemberInput
	var enqueue enqueueOutcome
	skipPlacement := false

	if partyInput != nil && len(partyInput.Members) > 0 {
		members = partyInput.Members
	} else {
		members = []JoinPartyMemberInput{{UserID: callerID, QueuePath: queuePath}}
	}

	if existing, err := getWaitingQueueEntryForUserTx(ctx, tx, modeQueueID, callerID); err == nil {
		if sameQueuePath(existing.QueuePath, queuePath) {
			skipPlacement = true
			enqueue.alreadyInQueue = true
		}
	} else if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	if !skipPlacement {
		for _, member := range members {
			if err := s.prepareCatalogPartyJoinTx(ctx, tx, member.UserID); err != nil {
				return nil, err
			}
		}
		if partyInput != nil && len(partyInput.Members) > 0 {
			tree := partyInput.Tree
			if len(tree.AllMembers()) == 0 {
				return nil, fmt.Errorf("store: party tree is required")
			}
			created, err := s.CreatePartyFromTreeTx(ctx, tx, modeQueueID, callerID, tree, members)
			if err != nil {
				return nil, err
			}
			partyInput.Tree = tree
			for _, member := range members {
				out, err := enqueueModeQueueTx(ctx, tx, joinCtx.Game.ID, modeQueueID, member.UserID, member.QueuePath, &created.ID)
				if err != nil {
					return nil, err
				}
				if member.UserID == callerID {
					enqueue = out
				}
			}
		} else {
			created, err := s.CreateSoloPartyTx(ctx, tx, modeQueueID, callerID, SoloPartyInput{
				QueuePath: queuePath,
			})
			if err != nil {
				return nil, err
			}
			enqueue, err = enqueueModeQueueTx(ctx, tx, joinCtx.Game.ID, modeQueueID, callerID, queuePath, &created.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	waiting, err := listWaitingModeQueueEntriesTx(ctx, tx, modeQueueID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &QueueJoinResult{
		GameID:         joinCtx.Game.ID,
		ModeQueueID:    modeQueueID,
		Status:         QueueStatusWaiting,
		QueuedCount:    len(waiting),
		QueuePath:      optionalQueuePathRef(queuePath),
		NotifyUserIDs:  notifyUserIDsFromMembers(members),
		AlreadyInQueue: enqueue.alreadyInQueue,
		SwitchedFrom:   enqueue.switchedFrom,
	}, nil
}

func notifyUserIDsFromMembers(members []JoinPartyMemberInput) []uuid.UUID {
	ids := make([]uuid.UUID, len(members))
	for i, m := range members {
		ids[i] = m.UserID
	}
	return ids
}

func (s *Store) markPartyStatusTx(ctx context.Context, tx *sql.Tx, partyID uuid.UUID, status string) error {
	_, err := tx.ExecContext(ctx, `UPDATE parties SET status = $2 WHERE id = $1`, partyID, status)
	return err
}

func (s *Store) linkPartyQueuesToFormingTx(ctx context.Context, tx *sql.Tx, partyID, formingMatchID uuid.UUID) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE game_queues
		SET forming_match_id = $1
		WHERE party_id = $2 AND status = 'waiting'
	`, formingMatchID, partyID)
	return err
}

func (s *Store) fireFormingMatchTx(
	ctx context.Context,
	tx *sql.Tx,
	joinCtx *ModeQueueJoinContext,
	fm *FormingMatch,
) (*formingFireResult, error) {
	assignments, err := s.ListFormingAssignmentsTx(ctx, tx, fm.ID)
	if err != nil {
		return nil, err
	}

	session, err := createModeQueueSessionTx(ctx, tx, joinCtx.Game.ID, joinCtx.Mode.ID, joinCtx.ModeQueue.ID)
	if err != nil {
		return nil, err
	}
	if err := completePriorActiveSessionsForModeQueueTx(ctx, tx, joinCtx.ModeQueue.ID, session.ID, time.Now()); err != nil {
		return nil, err
	}

	notifyIDs := make([]uuid.UUID, 0, len(assignments))
	tableIDs := make([]uuid.UUID, 0)
	seen := make(map[uuid.UUID]struct{})
	tableSeen := make(map[uuid.UUID]struct{})
	partySeen := make(map[uuid.UUID]struct{})

	for _, assignment := range assignments {
		if assignment.UserID == nil {
			continue
		}
		userID := *assignment.UserID
		entry, err := getWaitingQueueEntryForUserTx(ctx, tx, joinCtx.ModeQueue.ID, userID)
		if err != nil {
			return nil, err
		}
		if err := markQueueEntryMatchedTx(ctx, tx, entry.ID); err != nil {
			return nil, err
		}
		returnCtx := CatalogLFGReturnContext(joinCtx.Game.ID, joinCtx.ModeQueue.ID)
		if err := addSessionParticipantTx(ctx, tx, session.ID, userID, assignment.SeatKey, returnCtx); err != nil {
			return nil, err
		}
		if _, ok := seen[userID]; !ok {
			seen[userID] = struct{}{}
			notifyIDs = append(notifyIDs, userID)
		}
		if assignment.PartyID != nil {
			if _, ok := partySeen[*assignment.PartyID]; !ok {
				partySeen[*assignment.PartyID] = struct{}{}
				if err := s.markPartyStatusTx(ctx, tx, *assignment.PartyID, PartyStatusMatched); err != nil {
					return nil, err
				}
			}
		}
		if assignment.TableID != nil {
			if _, ok := tableSeen[*assignment.TableID]; !ok {
				tableSeen[*assignment.TableID] = struct{}{}
				tableIDs = append(tableIDs, *assignment.TableID)
			}
		}
	}

	if err := s.MarkFormingMatchFiredTx(ctx, tx, fm.ID, time.Now()); err != nil {
		return nil, err
	}

	return &formingFireResult{session: session, notifyIDs: notifyIDs, tableIDs: tableIDs}, nil
}

func getWaitingQueueEntryForUserTx(ctx context.Context, tx *sql.Tx, modeQueueID, userID uuid.UUID) (*QueueEntry, error) {
	row := tx.QueryRowContext(ctx, `
		SELECT `+queueColumns+`
		FROM game_queues
		WHERE mode_queue_id = $1 AND user_id = $2 AND status = 'waiting'
		ORDER BY joined_at DESC
		LIMIT 1
	`, modeQueueID, userID)
	return scanQueueEntry(row)
}
