# Pub/Sub and queue notifications

JoinQuest uses Redis pub/sub so every API instance can deliver queue match events to subscribed clients (part of the shared lobby experience — see [vision](./vision.md)).

## Channels

Per-user queue updates:

```text
lobby:user:{userId}:queue
```

Payload JSON (`internal/pubsub/events.go`):

- `queueId`, `status` (`WAITING` | `MATCHED` | `LEFT`)
- `sessionId`, `joinUrl` when matched
- `queuedCount` while waiting

## Configuration

| Variable | Purpose |
|----------|---------|
| `REDIS_URL` | e.g. `redis://127.0.0.1:6379/0` — omit to use in-memory broker (single process only) |
| `GAME_CLIENT_BASE_URL` | Fallback browser `play_url` when a game row has no `play_url` |
| `LOBBY_GAME_TOKEN_PEPPER` | HMAC key for per-game `lobby.serviceToken` (`v1.{gameId}.{sig}`) on provision and `registerGame` |
| `MAGIC_LINK_PEPPER` | Optional pepper for hashed magic-link tokens at rest (recommended in production) |
| `LOBBY_GAME_SERVICE_TOKEN` | Dev-only global Bearer fallback when pepper is unset |
| `LOBBY_ISSUER_URL` | Canonical issuer URL for seat JWT `iss` and provision `lobbyId` (falls back to `LOBBY_PUBLIC_URL`, then `http://localhost:8080`) |

Match launch URLs use the protocol in [`lobby-protocol-handoff.md`](./lobby-protocol-handoff.md):
`{playUrl}?match=<sessionId>&token=<seat-jwt>`.

**Provision ownership:** only the forming worker (and table `startTable`) call `POST /api/v1/matches`. GraphQL queries/subscriptions read stored launch URL bases — they never trigger provision. On transient game-server errors the worker retries with backoff (100ms → 500ms → 2s → 5s → 15s → 30s) without rolling back the matched session.

Local `./scripts/dev.sh` starts Postgres and Redis via Docker and sets both variables.

## Debugging queue delivery

Trace the full path from match fire → Redis → GraphQL subscription → browser.

### Backend (`LOBBY_PUBSUB_DEBUG=true`)

Logs to backend stdout:

| Log prefix | Meaning |
|------------|---------|
| `pubsub: reconcile` | Forming worker evaluated a queue (fired or not) |
| `pubsub: forming matched` | Match provisioned; about to publish per-user events |
| `pubsub: publish` | Event written to a user's Redis channel |
| `pubsub: redis publish/receive` | Redis broker I/O |
| `pubsub: subscription open/initial/deliver` | GraphQL `queueUpdated` resolver lifecycle |

```bash
# local
LOBBY_PUBSUB_DEBUG=true ./scripts/dev.sh

# joinquest
kubectl set env deployment/lobby-backend -n joinquest LOBBY_PUBSUB_DEBUG=true --containers=backend
kubectl logs -n joinquest deployment/lobby-backend -f | grep 'pubsub:'
```

### Frontend (browser console)

Enable any one of:

- `REACT_APP_LOBBY_DEBUG=true` in frontend env (see `k8s/env/joinquest.yaml`)
- `?lobbyDebug=1` on the URL (handy on iOS)
- `localStorage.setItem('lobbyDebug', '1')` then reload

Console lines are prefixed `[lobby:…]`:

| Tag | Meaning |
|-----|---------|
| `queue:ws:connected` | WebSocket to `/graphql` is up |
| `queue:subscribe:update` | `queueUpdated` payload received |
| `intent:queue:ws-update` | Banner hook handling a WS update |
| `intent:refresh:done` with `reason: poll-waiting` | HTTP poll picked up state (fallback path) |
| `intent:refresh:done` with `reason: ws-followup` | Refresh after a WS event |

**How to read a two-player test:** you should see backend `publish` + `subscription deliver` within ~50ms of the second join, and frontend `queue:subscribe:update` at roughly the same time. If backend delivers but the browser only logs `poll-waiting`, the WebSocket path is broken. If backend never logs `publish`, the worker or provision step failed.

## GraphQL

- **Mutation:** `joinQueue(queueId)` — enqueue on a catalog mode queue; match when `playersToStart` met
- **Mutation:** `leaveQueue(queueId)`
- **Subscription:** `queueUpdated(queueId)` — WebSocket via `/graphql`; send `Authorization: Bearer <jwt>` in `connection_init` (see `subscriptionAuth` query)

## Game projects

Reuse the same pattern: Redis channel per room or user, JSON events, subscribe from Node with `ioredis` or `redis` package. Keep gameplay state in the game DB; use pub/sub only for notifications and cross-service fan-out.

**Room chat (JoinQuest):** channel `lobby:room:{roomId}` — payload `RoomEvent` with `type` `updated` (membership) or `message` (new chat row). GraphQL subscriptions: `roomUpdated`, `roomMessageAdded`.
