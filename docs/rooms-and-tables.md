# Rooms and tables

Social **rooms** for friends to gather; **tables** are forming private games inside a room.

**Related:** [seat-templates-and-matchmaking.md](./seat-templates-and-matchmaking.md) · [game-catalog-architecture.md](./game-catalog-architecture.md) · [player-return-routing.md](./player-return-routing.md) · [composition-and-join-options.md](./composition-and-join-options.md)

---

## Roadmap

| Step | Scope |
|------|--------|
| **Step 1 — Rooms** (shipped) | Chat rooms: create/join, short invite code, member list, **QR** + easy link sharing. One room per player. |
| **Step 2 — Tables** (shipped) | 0..N forming tables per room; **seat-level sitting** (`sitAtTable(tableId, seatKey)`); king controls; queue ↔ table mutual exclusion; per-mode catalog actions; lazy stale discard. |
| **Step 3 — Table LFG** (shipped) | King **Look for group** backfill; `formingGaps` on TableCard; catalog fills remaining roles. |
| **Step 4+** | Phase C: affinity-aware weighted dequeue. |

---

## Step 1: Rooms (chat + invite)

See git history / prior docs for Step 1 detail. Rooms remain the social shell: invite code, chat, share/QR, one room per player globally.

**Step 2 change:** sitting at a table or joining a queue are **mutually exclusive** (see below). Chat-only room membership still works alongside either intent.

---

## Step 2: Tables (forming games)

### Concepts

| Entity | Purpose |
|--------|---------|
| **Table** | One **forming** match for a **game + mode** inside a room. Many tables per room (unbounded). |
| **Seat** | A specific expanded `seatKey` from the mode template (not a queue path alone). |
| **King** | Longest-sitting player (`min(seated_at)`). Recomputed when anyone leaves. Shown in UI (“King: Pat”). Controls **Start now** and **Discard** when eligible. |

```text
Room (invite code, chat)
  └── Table × N (forming)
        ├── game + mode
        ├── seat-level sitting (seatKey)
        ├── king: Start now; **Look for group** backfill (Phase B)
        └── on Start → session + provision (return to room)
```

### Player intent (global)

| Slot | Rule |
|------|------|
| **Table seat** | One active table seat globally. Sitting **leaves** any other table seat. |
| **Queue vs table** | Mutually exclusive — `sitAtTable` clears waiting queue; `joinQueue` clears table seat. |
| **Matched queue** | Cannot sit at a table until the active match is left. |

### Sitting and caps

- Join via **`sitAtTable(tableId, seatKey)`** — validates the key exists on the mode, seat is empty, path `max` not exceeded, and total ≤ mode `maxPlayers`.
- **Start** requires king, total ∈ `[minPlayers, maxPlayers]`, and every template path meets its **minimum** (path `min`, else `playersToStart` for that path).
- **Team layout in UI** uses `seatKey` prefix (e.g. `Team-1` vs `Team-2`). Role caps come from shared queue paths (e.g. Overwatch `DPS` cap 4 across both teams). **`affinityKey` is not used in Step 2.**

### Catalog UX

Per **mode** row under each game:

- **Look for group** — existing queue join (unchanged).
- **Create private game** — `createPrivateTable(gameId, modeId)` creates a room if needed, sweeps stale empty tables, adds a forming table.

### Table UI (room panel)

- **Tables (N)** section above chat.
- **TableCard:** game/mode, king badge, team columns when `Team-N` prefixes appear, individual seat chips, king-only **Start now**, **Look for group** backfill (when seated), role **formingGaps**, **Discard** when eligible.
- **Intent banner:** `myActiveIntent` (catalog wait / matched) or `myTableSeat` (table forming / backfill).

### Stale tables

- Discardable when **zero seated** and (**≥ 60s old** OR king manually discards).
- Auto-sweep on **create table** removes stale empty tables first.
- No countdown UI; button label **Discard**.

### Return routing

```json
{
  "kind": "room",
  "path": "/room/X7K2M9",
  "gameId": "…",
  "roomId": "…",
  "tableId": "…"
}
```

### Data model

#### `room_tables`

| Column | Notes |
|--------|--------|
| `id`, `room_id`, `game_id`, `mode_id` | Many per room |
| `status` | `forming`, `started`, `discarded` |
| `session_id` | Set on start |

#### `table_seats`

| Column | Notes |
|--------|--------|
| `table_id`, `user_id` | **Unique `user_id`** globally |
| `seat_key` | Expanded template key |
| `seated_at` | King tie-break |

### GraphQL (Step 2)

```graphql
type Table {
  id: ID!
  game: Game!
  mode: GameMode!
  createdAt: Time!
  king: User
  seats: [TableSeat!]!
  seatSlots: [TableSeatSlot!]!
  canStart: Boolean!
  canDiscard: Boolean!
  lookForGroupOptions: [TableLookForGroupOption!]!
}

extend type Room { tables: [Table!]! }

extend type Query { myTableSeat: MyTableSeat }

extend type Mutation {
  createPrivateTable(gameId: ID!, modeId: ID!): Table!
  createTable(roomId: ID!, gameId: ID!, modeId: ID!): Table!
  sitAtTable(tableId: ID!, seatKey: String!): Table!
  leaveTable(tableId: ID!): Boolean!
  discardTable(tableId: ID!): Boolean!
  startTable(tableId: ID!): JoinResult!
}

extend type Subscription {
  tableUpdated(roomId: ID!): Table!
}
```

Realtime: Redis `lobby:room:{roomId}` publishes `table_updated` events (same channel as chat/membership).

### Template validation

Catalog sync rejects `seatTemplate` when two queue paths share the same **`displayName`** (case-insensitive). Paths are **not** merged in the UI.

---

## Locked product defaults

| Decision | Choice |
|----------|--------|
| Build order | Rooms first, tables second |
| Tables per room | 0..N |
| Sitting | **Seat-level** (`seatKey`) |
| King | Longest-sitting player; visible |
| Tables per player | 1 seat globally |
| Queue vs table | Mutually exclusive |
| Private game | No queue; binds game + mode only |
| Table LFG | **Look for group** backfill (Phase B) |
| Stale empty tables | Lazy sweep + manual Discard |
| Catalog | Per-mode LFG + Create private game |
