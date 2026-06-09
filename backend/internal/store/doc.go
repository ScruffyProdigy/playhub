// Package store provides typed PostgreSQL persistence for JoinQuest.
//
// File layout (by domain):
//
//   - Auth & users: user.go, magic_link.go
//   - Catalog: catalog.go, game.go, mode_seats.go
//   - Queues & join: queue.go, queue_tx.go, queue_status.go, queue_active.go,
//     queue_path_matchmaking.go, mode_queue_matchmaking.go, queue_expire.go
//   - LFG forming (Phase B): forming_match.go, forming_join.go, forming_place.go,
//     forming_gaps.go, party.go, table_backfill.go
//   - Tables & rooms: table.go, room.go, active_game.go
//   - Sessions: session.go, session_participants.go, match_lifecycle.go, match_rollback.go
//   - Goods: digital_good.go, inventory.go
//   - Models: models_*.go
package store
