# Reference games

Working JoinQuest integrations you can play, read, or fork. Use them to learn the handshake **and** what a shippable multiplayer client looks like after integration.

| Game | GitHub | Play live | Good template for |
|------|--------|-----------|-------------------|
| **Rock Paper Scissors Lizard Robot** | [ScruffyProdigy/rpslr](https://github.com/ScruffyProdigy/rpslr) | [rpsls-duel.win](https://rpsls-duel.win) | Smallest duel game — API + client + integration tests; start here for the handshake |
| **Word Hunt** | [ScruffyProdigy/wordhunt](https://github.com/ScruffyProdigy/wordhunt) | [word-hunt-arena.win](https://word-hunt-arena.win) | Multi-player party game — path-style launch URLs, richer client UX, realtime sync |

Both repos include:

- Game API (`healthz`, `status`, `game-modes`, provision, JWT claim)
- A **playable web client** at the launch URLs JoinQuest opens after matchmaking
- Local tests mirroring [integration guide §8](./developer-integration-guide.md#8-recommended-local-tests-game-repo)
- Return-to-JoinQuest links after a match

## Who should use these

**Human developers** — skim a live game, then clone the repo closest to your idea. You do not have to fork unchanged; copy patterns for API layout, client boot from `?token=`, and lobby return.

**AI agents** — after Phase 5 (checks green) or in parallel with Phase 3, read the reference repo that best matches player count and session shape. Use it for integration test structure **and** for client/gameplay patterns in Phase 9 (see [developer-agent-playbook.md](./developer-agent-playbook.md)).

## Related

- [Developer integration guide](./developer-integration-guide.md)
- [Game-minted launch URLs](./game-minted-launch-urls.md)
- [Lobby ↔ game handoff](./lobby-protocol-handoff.md)
