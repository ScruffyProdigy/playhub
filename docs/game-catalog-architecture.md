# Game catalog, modes, and queues

Locked architecture decisions for **third-party game integration** on JoinQuest.

**Why this matters for game authors:** Players discover your title in our catalog and queue here; we fill seats and hand off to your `play_url`. You publish modes and a seat layout; we own matchmaking. Product story: [`vision.md`](./vision.md). Wire protocol: [`lobby-protocol-handoff.md`](./lobby-protocol-handoff.md).

## Model

```text
Game
  └── GameMode (many)
        ├── seatTemplate (from cached game-modes manifest; expanded to leaf seats)
        └── Queue (many)
              ├── single path (`""`): one “Look for group” bucket (e.g. 1v1 duel)
              └── composition: one bucket per `queuePath` (DPS, Tank, …)
```

- **Game** — `slug`, `play_url`, `api_base_url`, sync metadata.
- **GameMode** — one way to play; identified by `mode_key` sent on provision.
  A two-player mode is a mode with two seats, not a special “duel” type.
- **Queue** — matchmaking bucket scoped to `queue_id` (fifo default or role-specific).
- **Session** — `game_id`, `mode_id`, optional `queue_id` (null for table/room starts).
- **Room / Table** — social room (invite code) + forming table per game/mode; see [rooms-and-tables.md](./rooms-and-tables.md).

**Full spec (game authors + Lobby LFG):** [`seat-templates-and-matchmaking.md`](./seat-templates-and-matchmaking.md) —
`seatTemplate` → seat map (game contract); constraint-based LFG fills the map before provision.

## Seat fill (current implementation vs target)

**Today (shipped):** games publish **`seatTemplate` only** (flat `seats[]` is rejected on sync).
Lobby expands the template to leaf seats (`seat_key`, `queue_path`, `sort_order`), stores
`seat_template` on the mode, and auto-creates one default queue per mode. Matchmaking pairs
waiting players to seats by **`queue_path`** (fifo within each path). Modes with a single
path (often `""`) show one **Look for group** button; modes with multiple paths show **Join as …**
per role bucket. See [composition-and-join-options.md](./composition-and-join-options.md).

**Target (Phase B):** **forming match** filled by constraint-based LFG (parties, affinity
gaps, weighted dequeue, variable `sizeForQueue`); games still receive final `seatKey`
assignments only. See the linked doc.

## Admin registration

- Admin calls `registerGame` (email allowlist via `LOBBY_ADMIN_EMAILS`).
- Lobby immediately pulls `{apiBaseUrl}/healthz`, `/api/v1/status`, `/api/v1/game-modes`.
- **Fail closed** — do not list the game if health or game-modes fetch fails.
- Cache normalized modes/seats plus `manifest_hash`, `manifest_etag`, `manifest_synced_at`.
- Issue a per-game **webhook secret** for manifest-change callbacks.

## Manifest updates

- **Primary:** game `POST`s lobby `…/games/{slug}/manifest-changed` (authenticated with webhook secret).
- **Backup:** periodic poll of `/api/v1/game-modes` (ETag when supported).
- **In-flight matches** — never changed by a manifest refresh.
- **Added mode** — list in catalog; auto-create one default queue (`players_to_start = len(seats)`).
- **Removed mode** — disable mode and its queues.
- **Removed queue** — **kick all waiting players** with a user-visible message (e.g. via
  `queueUpdated` with status `LEFT` and reason text); do not start new matches on that queue.
- **Seat key changes** that break existing queues — disable affected queues and kick waiters with a message until an admin re-enables or recreates queues.

## Queues

- **Hybrid creation:** on register, auto-create queue(s) from template (default fifo, or
  one row per `queue: true` role dimension).
- Matchmaking is **scoped to `queue_id`**, not `game_id`.
- **One waiting queue per player** globally (migration `000012`). Joining another queue **automatically leaves** the previous waiting queue; the client shows a `message` on `joinQueue`.
- **Fifo:** FIFO within one queue until `players_to_start` (default `max`).
- **Composition:** match fires when every role bucket reaches `required`; coordinator
  assigns teams (see seat-templates doc).

## URLs

- **`play_url` and `api_base_url` live on the game** (shared across modes) until a title needs per-mode URLs.
- **`api_base_url`**: Game API **origin** only (e.g. `http://localhost:3001`, `https://rpsls-duel.win`) — same value as JWT `aud`. Lobby posts to `{api_base_url}/api/v1/matches`. Do not include the `/api` ingress prefix here.

## Provision

- `lobbyId` = `LobbyIssuer()`; `lobby.returnUrl` = browser Lobby URL; `lobby.graphqlUrl` = `{issuer}/graphql`; `lobby.serviceToken` = per-game credential `v1.{gameUUID}.{hmac}` from `LOBBY_GAME_TOKEN_PEPPER` (also returned by admin `registerGame.serviceToken`). Games store it per match; provision POST uses `Authorization: Bearer` with the same value. Legacy global `LOBBY_GAME_SERVICE_TOKEN` is dev-only fallback.
- `gameMode` = mode’s `mode_key` (game server resolves mode-specific rules locally).
- Pass `team` / `role` on seats from the cached template when applicable.
- Player display names via `player(id)` at `lobby.graphqlUrl` using `lobby.serviceToken` from provision.
- Banlist handshake (`403` + `bannedLobbyUserIds`) unchanged.

## Player API (target)

```graphql
games {
  modes {
    composition { slotKind required waitingCount }
    queues { id name slotKind requiredTotal playersToStart waitingCount }
    seats { seatKey affinityKey slotKind }
  }
}
joinQueue(queueId: ID!, party: PartyInput): JoinResult!
queueUpdated(queueId: ID!): QueueUpdate!  # include message/reason when kicked
```

Details in [`seat-templates-and-matchmaking.md`](./seat-templates-and-matchmaking.md).

Non-queue start paths (party, direct) are **v2**; schema keeps `queue_id` nullable on sessions.

## Admin auth

- `LOBBY_ADMIN_EMAILS` comma-separated allowlist; signed-in session required.
- `me { isAdmin }` for client gating; phone demos use normal sign-in, not API keys.

## Environments

- One catalog row per environment (e.g. joinquest.cc vs local), each with its own URLs and cached manifest.

## Deferred

- Signed provision requests (allowlist `lobbyId` + shared Bearer today).
- Webhook when manifest changes.
- Populate `team` / `role` on provision seats from manifest.
- Seat map LFG engine (`seat-templates-and-matchmaking.md`).
