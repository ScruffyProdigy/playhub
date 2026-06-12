# JoinQuest developer integration guide

This guide is a stub for Phase A. The full checklist (provision, JWT, launch URLs) ships in Phase B.

## Quick start

1. **Health** — `GET {apiBaseUrl}/healthz` must return `ok`.
2. **Status** — `GET {apiBaseUrl}/api/v1/status` returns your game name and version.
3. **Game modes** — `GET {apiBaseUrl}/api/v1/game-modes` returns modes with valid `seatTemplate` JSON.

## Seat manifest

Each mode needs `minPlayers`, `maxPlayers`, and a `seatTemplate` Lobby can expand into join buckets.

## Testing with friends

When your dashboard checklist shows manifest checks passing, use **Create test table** on your developer dashboard. That opens a private room table — invite friends with your room link. Your game stays off the public catalog until you request public release.

## Need help?

Re-run checks on your developer dashboard after fixing issues. Friendly errors link back to sections here as we expand this guide.
