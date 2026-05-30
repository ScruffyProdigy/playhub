# Pub/Sub and queue notifications

Lobby uses Redis pub/sub so every API instance can deliver queue match events to subscribed clients.

## Channels

Per-user queue updates:

```text
lobby:user:{userId}:queue
```

Payload JSON (`internal/pubsub/events.go`):

- `gameId`, `status` (`WAITING` | `MATCHED` | `LEFT`)
- `sessionId`, `joinUrl` when matched
- `queuedCount` while waiting

## Configuration

| Variable | Purpose |
|----------|---------|
| `REDIS_URL` | e.g. `redis://127.0.0.1:6379/0` — omit to use in-memory broker (single process only) |
| `GAME_CLIENT_BASE_URL` | Fallback browser `play_url` when a game row has no `play_url` |
| `LOBBY_GAME_SERVICE_TOKEN` | Optional Bearer for game `POST /api/v1/matches` |

Match launch URLs use the protocol in [`lobby-protocol-handoff.md`](./lobby-protocol-handoff.md):
`{playUrl}?match=<sessionId>&token=<seat-jwt>`.

Local `./scripts/dev.sh` starts Postgres and Redis via Docker and sets both variables.

## GraphQL

- **Mutation:** `joinGame(gameId)` — enqueue + match when `min_players` met
- **Mutation:** `leaveQueue(gameId)`
- **Subscription:** `queueUpdated(gameId)` — WebSocket via `/graphql`; send `Authorization: Bearer <jwt>` in `connection_init` (see `subscriptionAuth` query)

## Game projects

Reuse the same pattern: Redis channel per room or user, JSON events, subscribe from Node with `ioredis` or `redis` package. Keep gameplay state in the game DB; use pub/sub only for notifications and cross-service fan-out.
