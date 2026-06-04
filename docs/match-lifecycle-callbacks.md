# Match lifecycle callbacks (game → Lobby)

Games run on their own origin after provision. Lobby still needs **lifecycle signals** so the player shell, stats, and queue rules stay correct.

**Related:** provision + `returnUrl` in [lobby-protocol-handoff.md](./lobby-protocol-handoff.md).

---

## Two levels of reporting

| Callback | When | Purpose |
|----------|------|---------|
| **`reportPlayerFinished`** | One player is **done** with the match for themselves | Eliminated early, finished first in a race while others play, forfeited, etc. |
| **`reportMatchResult`** | The **match** is over for everyone | Winner, scores, abandoned, cancelled |

A race might call `reportPlayerFinished` for 1st, 2nd, 3rd as each crosses the line, then `reportMatchResult` when the last state is resolved.

A battle royale might call `reportPlayerFinished` on each elimination and `reportMatchResult` when one player remains (or time expires).

A synchronous 1v1 might only call `reportMatchResult` when both sides are done.

---

## Why both matter

- **Player UX** — JoinQuest can show “You placed 2nd” and enable re-queue while others are still racing.
- **Queue rules** — A player who is **finished** should not block “one queue at a time” as if they were still in an active match; Lobby can clear or downgrade their `matched` row when the game reports finish.
- **Integrity** — Final `reportMatchResult` is the authoritative outcome for leaderboards and disputes.

---

## Proposed GraphQL (game server, Bearer `serviceToken`)

```graphql
enum PlayerFinishReason {
  COMPLETED
  ELIMINATED
  FORFEIT
  DISCONNECT
}

enum MatchResultStatus {
  COMPLETED
  CANCELLED
  ABANDONED
}

mutation reportPlayerFinished(
  matchId: ID!           # Lobby externalMatchId
  lobbyUserId: ID!
  reason: PlayerFinishReason!
  placement: Int         # optional, e.g. 1 = first
  metadata: JSON
): Boolean!

mutation reportMatchResult(
  matchId: ID!
  status: MatchResultStatus!
  winnerLobbyUserIds: [ID!]
  metadata: JSON
): Boolean!
```

Lobby validates `matchId` + `serviceToken` game id, idempotency via `jti` or monotonic event ids (TBD in implementation).

---

## Lobby behavior (target)

1. **`reportPlayerFinished`** — mark participant finished; if all players finished or match ended, transition session; clear user’s `matched` queue row so they can join another queue.
2. **`reportMatchResult`** — set session `ended`, publish any subscriber updates, optional stats pipeline.

Until these exist, games should still use **`returnUrl`** so players can navigate back manually; queue expiry handles stuck `matched` rows.

---

## Implementation status

**Planned** — not yet in the GraphQL schema. Listed in [development.md](./development.md) and [api.md](./api.md).

Recommended order: **`reportPlayerFinished`** first (unblocks re-queue + race UX), then **`reportMatchResult`** (authoritative close).
