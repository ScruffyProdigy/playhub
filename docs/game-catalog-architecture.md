# Game catalog, modes, and queues

Locked architecture decisions for third-party game integration. Wire protocol details
remain in [`lobby-protocol-handoff.md`](./lobby-protocol-handoff.md).

## Model

```text
Game
  └── GameMode (many)
        ├── seat template (from cached game-modes manifest)
        └── Queue (many)
              └── players_to_start (fixed N before a match starts)
```

- **Game** — `slug`, `play_url`, `api_base_url`, sync metadata.
- **GameMode** — one way to play; identified by `mode_key` sent on provision.
  A two-player mode is a mode with two seats, not a special “duel” type.
- **Queue** — one matchmaking bucket for a mode; starts a match when
  `players_to_start` distinct players are waiting.
- **Session** — `game_id`, `mode_id`, optional `queue_id` (null for non-queue starts later).

## Seat fill

- **Mode owns the seat template** (from `GET {apiBaseUrl}/api/v1/game-modes`, cached in lobby).
- **Queue sets `players_to_start` only** — assign the first N seats from the mode template in order.
- Queue-specific seat subsets are out of scope until a game needs them.

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

- **Hybrid creation:** on register, auto-create one default queue per mode; admin may add named queues later.
- Matchmaking is **scoped to `queue_id`**, not `game_id`.
- v1 matchmaking: **FIFO** within a queue; one waiting entry per user per queue.

## URLs

- **`play_url` and `api_base_url` live on the game** (shared across modes) until a title needs per-mode URLs.

## Provision

- `lobbyId` = `LobbyIssuer()`; `lobby.returnUrl` = browser Lobby URL; `lobby.graphqlUrl` = `{issuer}/graphql`.
- `gameMode` = mode’s `mode_key` (game server resolves mode-specific rules locally).
- Pass `team` / `role` on seats from the cached template when applicable.
- Player display names via `player(id)` at `lobby.graphqlUrl`; match results via GraphQL when implemented (same service token).
- Banlist handshake (`403` + `bannedLobbyUserIds`) unchanged.

## Player API (target)

```graphql
games { modes { queues { id name playersToStart waitingCount } } }
joinQueue(queueId: ID!): JoinResult!
queueUpdated(queueId: ID!): QueueUpdate!  # include message/reason when kicked
```

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
