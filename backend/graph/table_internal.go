package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/graph/model"
	"github.com/scruffyprodigy/playhub/internal/pubsub"
	"github.com/scruffyprodigy/playhub/internal/seattemplate"
	"github.com/scruffyprodigy/playhub/internal/store"
)

func toGraphQLTable(table *store.RoomTable) *model.Table {
	if table == nil {
		return nil
	}
	return &model.Table{
		ID:        table.ID.String(),
		CreatedAt: table.CreatedAt,
	}
}

func (r *Resolver) loadTableModel(ctx context.Context, tableID uuid.UUID) (*model.Table, error) {
	st, err := r.requireStore()
	if err != nil {
		return nil, err
	}
	table, err := st.GetRoomTableByID(ctx, tableID)
	if err != nil {
		return nil, err
	}
	return toGraphQLTable(table), nil
}

func (r *Resolver) publishTableUpdated(ctx context.Context, roomID, tableID uuid.UUID) error {
	if r.PubSub == nil {
		return nil
	}
	return pubsub.PublishRoomEvent(ctx, r.PubSub, roomID.String(), pubsub.RoomEvent{
		Type:    pubsub.RoomEventTableUpdated,
		TableID: tableID.String(),
	})
}

func (r *Resolver) tableAfterMutation(ctx context.Context, table *store.RoomTable) (*model.Table, error) {
	if err := r.publishTableUpdated(ctx, table.RoomID, table.ID); err != nil {
		return nil, err
	}
	return toGraphQLTable(table), nil
}

func seatDisplayName(modeSeats []store.GameModeSeat, seatKey string, template json.RawMessage) string {
	for _, seat := range modeSeats {
		if seat.SeatKey != seatKey {
			continue
		}
		path := ""
		if seat.QueuePath != nil {
			path = strings.TrimSpace(*seat.QueuePath)
		}
		if path != "" {
			roleLabel := path
			if specs, err := seattemplate.PathSpecs(template); err == nil {
				for _, spec := range specs {
					if spec.QueuePath == path {
						roleLabel = spec.DisplayName
						break
					}
				}
			}
			if suffix := seatKeySuffix(seatKey, path); suffix != "" {
				return roleLabel + " · " + suffix
			}
			return roleLabel
		}
		if seat.Role != nil && strings.TrimSpace(*seat.Role) != "" {
			return strings.TrimSpace(*seat.Role)
		}
		return seatKey
	}
	return seatKey
}

func seatKeySuffix(seatKey, queuePath string) string {
	prefix := queuePath + "-"
	if strings.HasPrefix(seatKey, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(seatKey, prefix))
	}
	return ""
}

func (r *Resolver) buildTableSeatSlots(ctx context.Context, tableID uuid.UUID) ([]*model.TableSeatSlot, error) {
	st, err := r.requireStore()
	if err != nil {
		return nil, err
	}
	table, err := st.GetRoomTableByID(ctx, tableID)
	if err != nil {
		return nil, err
	}
	mode, err := st.GetGameModeByID(ctx, table.ModeID)
	if err != nil {
		return nil, err
	}
	modeSeats, err := st.ListGameModeSeats(ctx, table.ModeID)
	if err != nil {
		return nil, err
	}
	seated, err := st.ListTableSeats(ctx, tableID)
	if err != nil {
		return nil, err
	}
	byKey := map[string]uuid.UUID{}
	for _, s := range seated {
		byKey[s.SeatKey] = s.UserID
	}

	out := make([]*model.TableSeatSlot, len(modeSeats))
	for i, seat := range modeSeats {
		slot := &model.TableSeatSlot{
			SeatKey:     seat.SeatKey,
			DisplayName: seatDisplayName(modeSeats, seat.SeatKey, mode.SeatTemplate),
			TeamPrefix:  storeTeamPrefix(seat.SeatKey),
		}
		if seat.QueuePath != nil {
			slot.QueuePath = seat.QueuePath
		}
		if uid, ok := byKey[seat.SeatKey]; ok {
			user, uErr := st.GetUserByID(ctx, uid)
			if uErr == nil {
				slot.User = ToGraphQLUser(user)
			}
		}
		out[i] = slot
	}
	return out, nil
}

func storeTeamPrefix(seatKey string) *string {
	prefix := teamPrefixFromSeatKey(seatKey)
	if prefix == "" {
		return nil
	}
	return &prefix
}

func teamPrefixFromSeatKey(seatKey string) string {
	parts := strings.Split(seatKey, "-")
	if len(parts) >= 2 && strings.EqualFold(parts[0], "Team") {
		return parts[0] + "-" + parts[1]
	}
	return ""
}

func (r *Resolver) buildLookForGroupOptions(ctx context.Context, tableID uuid.UUID) ([]*model.TableLookForGroupOption, error) {
	st, err := r.requireStore()
	if err != nil {
		return nil, err
	}
	table, err := st.GetRoomTableByID(ctx, tableID)
	if err != nil {
		return nil, err
	}
	canStart, err := st.TableCanStart(ctx, tableID)
	if err != nil {
		return nil, err
	}
	seated, err := st.ListTableSeats(ctx, tableID)
	if err != nil {
		return nil, err
	}
	mode, err := st.GetGameModeByID(ctx, table.ModeID)
	if err != nil {
		return nil, err
	}
	visibleBase := len(seated) >= 1 && !canStart && len(seated) < mode.MaxPlayers

	queues, err := st.ListModeQueuesByModeID(ctx, table.ModeID)
	if err != nil {
		return nil, err
	}
	out := make([]*model.TableLookForGroupOption, 0, len(queues))
	for _, q := range queues {
		if q.Status != store.ModeQueueActive {
			continue
		}
		out = append(out, &model.TableLookForGroupOption{
			QueueID:   q.ID.String(),
			QueueName: q.Name,
			Visible:   visibleBase,
			Enabled:   false,
		})
	}
	return out, nil
}

func toGraphQLMyTableSeat(view *store.UserTableSeatView, seatDisplayName string) *model.MyTableSeat {
	if view == nil {
		return nil
	}
	return &model.MyTableSeat{
		TableID:         view.TableID.String(),
		RoomID:          view.RoomID.String(),
		InviteCode:      view.InviteCode,
		GameID:          view.GameID.String(),
		GameName:        view.GameName,
		ModeID:          view.ModeID.String(),
		ModeName:        view.ModeName,
		SeatKey:         view.SeatKey,
		SeatDisplayName: seatDisplayName,
	}
}

func (r *mutationResolver) startTableInternal(ctx context.Context, tableID uuid.UUID) (*model.JoinResult, error) {
	st, err := r.requireStore()
	if err != nil {
		return nil, err
	}
	userID, err := requireAuthUserID(ctx)
	if err != nil {
		return nil, err
	}

	table, err := st.GetRoomTableByID(ctx, tableID)
	if err != nil {
		return nil, err
	}

	result, err := st.StartTable(ctx, tableID, userID)
	if err != nil {
		return nil, err
	}

	game, err := st.GetGameByID(ctx, result.GameID)
	if err != nil {
		return nil, err
	}
	launchURLs, err := r.finalizeMatchedSession(ctx, game, result.SessionID, result.NotifyUserIDs)
	if err != nil {
		return nil, err
	}
	if err := r.publishTableUpdated(ctx, table.RoomID, tableID); err != nil {
		return nil, err
	}

	launch := launchURLs[userID]
	sessionID := result.SessionID.String()
	return &model.JoinResult{
		Queued:    false,
		SessionID: &sessionID,
		JoinURL:   &launch,
	}, nil
}

func (r *Resolver) requireTableRoomMember(ctx context.Context, tableID, userID uuid.UUID) (*store.RoomTable, error) {
	st, err := r.requireStore()
	if err != nil {
		return nil, err
	}
	table, err := st.GetRoomTableByID(ctx, tableID)
	if err != nil {
		return nil, err
	}
	if err := r.requireRoomMember(ctx, table.RoomID, userID); err != nil {
		return nil, err
	}
	return table, nil
}

func tableRoomIDFromObj(obj *model.Table) (uuid.UUID, error) {
	if obj == nil || strings.TrimSpace(obj.ID) == "" {
		return uuid.Nil, fmt.Errorf("table is required")
	}
	tableID, err := parseUUID(obj.ID, "table id")
	if err != nil {
		return uuid.Nil, err
	}
	return tableID, nil
}
