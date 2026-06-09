package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/lfg/partytree"
)

// SoloPartyInput describes a single-player party for queue join.
type SoloPartyInput struct {
	QueuePath string
}

// JoinPartyMemberInput is one member in a queue join party.
type JoinPartyMemberInput struct {
	UserID    uuid.UUID
	QueuePath string
}

// JoinPartyInput describes a party joining a mode queue.
type JoinPartyInput struct {
	Tree    partytree.Node
	Members []JoinPartyMemberInput
	Pinned  []partytree.PinnedSeat
	TableID *uuid.UUID
}

// CreatePartyFromTreeTx creates a party with a layout tree and member rows.
func (s *Store) CreatePartyFromTreeTx(
	ctx context.Context,
	tx *sql.Tx,
	modeQueueID uuid.UUID,
	leaderID uuid.UUID,
	tree partytree.Node,
	members []JoinPartyMemberInput,
) (*Party, error) {
	if len(members) == 0 {
		return nil, fmt.Errorf("store: party requires at least one member")
	}
	for _, member := range members {
		active, err := s.getActivePartyForUserTx(ctx, tx, member.UserID)
		if err != nil && !errors.Is(err, ErrNotFound) {
			return nil, err
		}
		if active != nil {
			return nil, fmt.Errorf("store: user %s already in an active party", member.UserID)
		}
	}

	treeJSON, err := json.Marshal(tree)
	if err != nil {
		return nil, fmt.Errorf("store: encode party tree: %w", err)
	}

	row := tx.QueryRowContext(ctx, `
		INSERT INTO parties (leader_user_id, mode_queue_id, status, party_tree)
		VALUES ($1, $2, $3, $4)
		RETURNING id, leader_user_id, mode_queue_id, status, party_tree, created_at
	`, leaderID, modeQueueID, PartyStatusWaiting, treeJSON)
	party, err := scanPartyRow(row)
	if err != nil {
		return nil, err
	}

	for i, member := range members {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO party_members (party_id, user_id, queue_path, sort_order)
			VALUES ($1, $2, $3, $4)
		`, party.ID, member.UserID, strings.TrimSpace(member.QueuePath), i); err != nil {
			return nil, err
		}
	}
	return party, nil
}

// CreateSoloPartyTx creates a one-member party for queue join.
func (s *Store) CreateSoloPartyTx(
	ctx context.Context,
	tx *sql.Tx,
	modeQueueID, userID uuid.UUID,
	input SoloPartyInput,
) (*Party, error) {
	tree := partytree.SoloNode(userID.String(), input.QueuePath)
	return s.CreatePartyFromTreeTx(ctx, tx, modeQueueID, userID, tree, []JoinPartyMemberInput{{
		UserID:    userID,
		QueuePath: input.QueuePath,
	}})
}

func scanPartyRow(row interface{ Scan(dest ...any) error }) (*Party, error) {
	var party Party
	var leader sql.NullString
	var treeJSON []byte
	if err := row.Scan(
		&party.ID, &leader, &party.ModeQueueID, &party.Status, &treeJSON, &party.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if leader.Valid {
		id, err := uuid.Parse(leader.String)
		if err != nil {
			return nil, err
		}
		party.LeaderUserID = &id
	}
	if len(treeJSON) > 0 {
		party.PartyTree = treeJSON
	}
	return &party, nil
}

func (s *Store) getActivePartyForUserTx(ctx context.Context, tx *sql.Tx, userID uuid.UUID) (*Party, error) {
	row := tx.QueryRowContext(ctx, `
		SELECT p.id, p.leader_user_id, p.mode_queue_id, p.status, p.party_tree, p.created_at
		FROM parties p
		INNER JOIN party_members pm ON pm.party_id = p.id
		WHERE pm.user_id = $1 AND p.status IN ($2, $3)
		ORDER BY p.created_at DESC
		LIMIT 1
	`, userID, PartyStatusWaiting, PartyStatusPlaced)
	return scanPartyRow(row)
}

// CancelPartyTx marks a party cancelled.
func (s *Store) CancelPartyTx(ctx context.Context, tx *sql.Tx, partyID uuid.UUID) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE parties SET status = $2 WHERE id = $1 AND status IN ($3, $4)
	`, partyID, PartyStatusCancelled, PartyStatusWaiting, PartyStatusPlaced)
	if err != nil {
		return err
	}
	return ensureRowsAffected(result, ErrNotFound)
}

func decodePartyTree(raw json.RawMessage) (partytree.Node, error) {
	if len(raw) == 0 {
		return partytree.Node{}, nil
	}
	var tree partytree.Node
	if err := json.Unmarshal(raw, &tree); err != nil {
		return partytree.Node{}, fmt.Errorf("store: decode party tree: %w", err)
	}
	return tree, nil
}