# Rooms and tables

Social **rooms** for friends to gather; **tables** (forming games) come in a later step.

**Related:** [seat-templates-and-matchmaking.md](./seat-templates-and-matchmaking.md) · [game-catalog-architecture.md](./game-catalog-architecture.md) · [player-return-routing.md](./player-return-routing.md) · [composition-and-join-options.md](./composition-and-join-options.md)

---

## Roadmap

| Step | Scope |
|------|--------|
| **Step 1 — Rooms** (shipped) | Chat rooms: create/join, short invite code, member list, **QR** + easy link sharing (copy, SMS, native share). One room per player. **No tables.** |
| **Step 2 — Tables** | Forming match per game+mode inside a room; bucket caps; host Start; queue ↔ table exclusion; disabled table LFG until Phase B. |
| **Step 3+** | Multiple tables per room; table LFG backfill; affinity-aware buckets. |

---

## Step 1: Rooms (chat + invite)

### Purpose

A lightweight place for friends to coordinate before (or while) using the catalog — “get in the same room, talk, share the link.” Game pick and match start move to **Step 2 (tables)**.

### Concepts

| Entity | Step 1 |
|--------|--------|
| **Room** | Persistent chat space. Short **invite code** → `/room/X7K2M9`. |
| **Table** | *Not shipped yet.* Schema/API should allow a room to gain tables later without a breaking migration. |

```text
Room (invite code, chat, share/QR)
  └── (Step 2) Table → game + mode forming → Start → session
```

### Membership

| Rule | Step 1 |
|------|--------|
| One room per player | Joining or creating a room **leaves** any previous room. |
| Queue while in room | **Allowed.** Room chat does not conflict with catalog queue membership (table/queue exclusion arrives with Step 2). |

### Invite and sharing

| Feature | Behavior |
|---------|----------|
| **Short code** | ~6 alphanumeric characters; unique; case-insensitive lookup. |
| **Join URL** | `{LOBBY_PUBLIC_URL}/room/{code}` — returned as `joinUrl` on `Room`. |
| **Copy link** | One-click copy to clipboard. |
| **QR code** | Button opens modal with QR encoding `joinUrl` (client-generated). |
| **Share sheet** | `navigator.share({ url, title })` when available (mobile / supported browsers). |
| **Text message** | `sms:?body=` with pre-filled invite text + URL (best-effort; platform-dependent). |

Suggested default share text: *“Join my room on JoinQuest: {joinUrl}”*

### Chat

| Piece | Approach |
|-------|----------|
| **Persistence** | `room_messages` table (room_id, user_id, body, created_at); paginated history on load. |
| **Realtime** | Redis pub/sub channel `lobby:room:{roomId}:chat` + GraphQL subscription `roomMessageAdded(roomId)`. Same pattern as [pubsub.md](./pubsub.md) queue events. |
| **Send** | Mutation `sendRoomMessage(roomId, body)` — auth required, must be room member. |
| **Limits** | Reasonable max body length (e.g. 2k chars); rate limit TBD. |

### Data model (Step 1)

#### `rooms`

| Column | Notes |
|--------|--------|
| `id` | UUID |
| `invite_code` | Short unique code |
| `host_user_id` | Creator |
| `status` | `open`, `closed` |
| `created_at`, `updated_at` | |

#### `room_members`

| Column | Notes |
|--------|--------|
| `room_id`, `user_id` | **Unique `user_id`** globally — one room per player |
| `joined_at` | |

#### `room_messages`

| Column | Notes |
|--------|--------|
| `id` | UUID |
| `room_id`, `user_id` | |
| `body` | Text |
| `created_at` | |

Index: `room_messages(room_id, created_at DESC)` for history.

*No `tables` / `table_seats` until Step 2.*

### GraphQL (Step 1)

```graphql
type Room {
  id: ID!
  inviteCode: String!
  joinUrl: String!
  host: User!
  members: [User!]!
  messages(limit: Int = 50, before: ID): [RoomMessage!]!
}

type RoomMessage {
  id: ID!
  author: User!
  body: String!
  createdAt: Time!
}

extend type Mutation {
  createRoom: Room!
  joinRoom(inviteCode: String!): Room!
  leaveRoom: Boolean!
  sendRoomMessage(roomId: ID!, body: String!): RoomMessage!
}

extend type Query {
  room(inviteCode: String!): Room
  myRoom: Room
}

extend type Subscription {
  roomUpdated(roomId: ID!): Room!      # member join/leave
  roomMessageAdded(roomId: ID!): RoomMessage!
}
```

`createRoom` takes **no game/mode** in Step 1 — pure social room.

### UI (Step 1)

| Surface | Content |
|---------|---------|
| **Catalog / nav** | “Create room” → creates room, navigates to `/room/{code}`. |
| **`/room/:code`** | Member list, chat transcript + composer, share toolbar (copy, QR, share, text). |
| **Deep link** | Unsigned visitor → sign-in → join room → land on room page. |

### Return routing (Step 1)

No table starts yet — return context `kind: "room"` is wired when **Step 2** starts matches from a table. Room page path remains `/room/{inviteCode}`.

---

## Step 2: Tables (forming games)

*Deferred — full spec preserved below for when Step 1 ships.*

### Concepts

| Entity | Purpose |
|--------|---------|
| **Table** | One forming match for a **game + mode** inside a room — seated players, bucket caps, host Start, (later) LFG backfill. |

```text
Room (invite code, chat)
  └── Table (Step 2 v1: exactly one per room)
        ├── game + mode
        ├── seated players (bucket per join slot)
        ├── host: Start; disabled Look for group until Phase B LFG
        └── on Start → session + provision
```

**Step 2 v1:** each room linked to **at most one** table. Do not assume 1:1 forever — multiple tables per room later.

### Player membership (Step 2 adds)

| Slot | Rule |
|------|------|
| **Table seat** | One active table seat globally. Sitting **leaves** any other table. |
| **Queue vs table** | Mutually exclusive — sit at table clears queue; `joinQueue` clears table seat. |

Room membership alone (chat only, not seated) still does not conflict with queue until the player sits.

### Join buckets and capacity

Table joins validated against mode **seat template** (`PathSpec` / `game_mode_seats`):

- Bucket key = `queuePath` (v1); team-specific caps via distinct paths (e.g. `Team2.DPS`) or Phase B `affinityKey`.
- **Reject sit** if `seated + 1 > path.max` or total > mode `maxPlayers`.
- **Reject start** unless host, total ∈ `[min, max]`, every bucket ∈ `[path.min, path.max]`, assignment succeeds.

### Table UI: Start vs Look for group

| Control | When |
|---------|------|
| **Start game** | Enabled when validation passes. |
| **Look for group** | **Disabled** until Phase B LFG. **Shown** when headcount ≤ LFG target and ≥1 bucket short vs `sizeForQueue`. |

### Return routing (Step 2)

```json
{
  "kind": "room",
  "path": "/room/X7K2M9",
  "gameId": "…",
  "roomId": "…",
  "tableId": "…"
}
```

### Data model (Step 2 adds)

#### `tables`

| Column | Notes |
|--------|--------|
| `id`, `room_id` | v1: unique `room_id` while one table per room |
| `game_id`, `mode_id` | |
| `host_user_id`, `status`, `target_size`, `session_id` | |

#### `table_seats`

| Column | Notes |
|--------|--------|
| `table_id`, `user_id` | Unique `user_id` globally |
| `queue_path`, `affinity_key` | Join bucket |

### GraphQL (Step 2 adds)

```graphql
type Table {
  id: ID!
  game: Game!
  mode: GameMode!
  host: User!
  status: TableStatus!
  seats: [TableSeat!]!
  bucketCounts: [TableBucketCount!]!
  canStart: Boolean!
  lookForGroupState: LookForGroupState!
}

extend type Room {
  tables: [Table!]!
  myTableSeat: TableSeat
}

extend type Mutation {
  createTable(roomId: ID!, gameId: ID!, modeId: ID!, targetSize: Int): Table!
  sitAtTable(tableId: ID!, queuePath: String, affinityKey: String): Table!
  leaveTable(tableId: ID!): Boolean!
  startTable(tableId: ID!): JoinResult!
}

extend type Subscription {
  tableUpdated(tableId: ID!): Table!
}
```

---

## Locked product defaults

| Decision | Choice |
|----------|--------|
| Build order | **Rooms first**, tables second |
| Invite format | Short code in `/room/{code}` |
| Rooms per player | 1 |
| Room vs queue (Step 1) | **No conflict** — can chat in room while queued |
| Tables per player (Step 2) | 1 seat |
| Queue vs table (Step 2) | Mutually exclusive |
| Share UX | Copy link, QR modal, native share, SMS body |
| Table start / caps / LFG | Step 2 + Phase B as above |
