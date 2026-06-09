package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/lfg"
	"github.com/scruffyprodigy/playhub/internal/lfg/partytree"
)

type formingFireResult struct {
	session   *Session
	notifyIDs []uuid.UUID
}

func (s *Store) joinModeQueueFormingTx(
	ctx context.Context,
	tx *sql.Tx,
	joinCtx *ModeQueueJoinContext,
	modeQueueID, callerID uuid.UUID,
	queuePath string,
	partyInput *JoinPartyInput,
) (*QueueJoinResult, error) {
	matchSeats, err := matchSeatsFromTemplate(joinCtx.Seats, joinCtx.Mode.SeatTemplate)
	if err != nil {
		return nil, err
	}
	if len(matchSeats) == 0 {
		matchSeats = joinCtx.Seats
	}

	var party *Party
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
		if partyInput != nil && len(partyInput.Members) > 0 {
			tree := partyInput.Tree
			if len(tree.AllMembers()) == 0 {
				return nil, fmt.Errorf("store: party tree is required")
			}
			party, err = s.CreatePartyFromTreeTx(ctx, tx, modeQueueID, callerID, tree, members)
			if err != nil {
				return nil, err
			}
			partyInput.Tree = tree
			for _, member := range members {
				out, err := enqueueModeQueueTx(ctx, tx, joinCtx.Game.ID, modeQueueID, member.UserID, member.QueuePath, &party.ID)
				if err != nil {
					return nil, err
				}
				if member.UserID == callerID {
					enqueue = out
				}
			}
		} else {
			party, err = s.CreateSoloPartyTx(ctx, tx, modeQueueID, callerID, SoloPartyInput{
				QueuePath: queuePath,
			})
			if err != nil {
				return nil, err
			}
			enqueue, err = enqueueModeQueueTx(ctx, tx, joinCtx.Game.ID, modeQueueID, callerID, queuePath, &party.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	fm, err := s.GetOrCreateFillingFormingMatchTx(ctx, tx, joinCtx, matchSeats)
	if err != nil {
		return nil, err
	}

	if !skipPlacement && party != nil {
		input := partyInput
		if input == nil {
			input = &JoinPartyInput{Tree: partytree.SoloNode(callerID.String(), queuePath), Members: members}
		}
		placed, err := s.tryPlacePartyOnFormingTx(ctx, tx, fm.ID, party, input)
		if err != nil {
			return nil, err
		}
		if placed {
			if err := s.markPartyStatusTx(ctx, tx, party.ID, PartyStatusPlaced); err != nil {
				return nil, err
			}
			if err := s.linkPartyQueuesToFormingTx(ctx, tx, party.ID, fm.ID); err != nil {
				return nil, err
			}
		}
	}

	gaps, err := s.FormingPathGapsTx(ctx, tx, fm)
	if err != nil {
		return nil, err
	}

	if lfg.ReadyToFire(gaps) {
		fired, err := s.fireFormingMatchTx(ctx, tx, joinCtx, fm)
		if err != nil {
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &QueueJoinResult{
			GameID:        joinCtx.Game.ID,
			ModeQueueID:   modeQueueID,
			Status:        QueueStatusMatched,
			SessionID:     &fired.session.ID,
			QueuedCount:   0,
			QueuePath:     optionalQueuePathRef(queuePath),
			NotifyUserIDs: fired.notifyIDs,
			SwitchedFrom:  enqueue.switchedFrom,
		}, nil
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
	seen := make(map[uuid.UUID]struct{})
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
	}

	if err := s.MarkFormingMatchFiredTx(ctx, tx, fm.ID, time.Now()); err != nil {
		return nil, err
	}

	return &formingFireResult{session: session, notifyIDs: notifyIDs}, nil
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
